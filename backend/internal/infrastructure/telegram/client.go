package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// Client — минимальный клиент Telegram Bot API (только нужные методы).
type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		baseURL: "https://api.telegram.org/bot" + token + "/",
		// Long polling держит запрос до 30 с, загрузка чека — до пары мегабайт.
		http: &http.Client{Timeout: 60 * time.Second},
	}
}

// --- Типы Bot API ---

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type Message struct {
	MessageID       int64  `json:"message_id"`
	MessageThreadID int64  `json:"message_thread_id"`
	Chat            Chat   `json:"chat"`
	From            *User  `json:"from"`
	Text            string `json:"text"`
	// В темах сообщение — ответ на служебное «тема создана», из него берём название темы
	ReplyToMessage    *Message    `json:"reply_to_message"`
	ForumTopicCreated *ForumTopic `json:"forum_topic_created"`
}

type ForumTopic struct {
	Name string `json:"name"`
}

// TopicName — название темы, в которой написано сообщение (если Telegram его прислал).
func (m *Message) TopicName() string {
	if m.ReplyToMessage != nil && m.ReplyToMessage.ForumTopicCreated != nil {
		return m.ReplyToMessage.ForumTopicCreated.Name
	}
	return ""
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    User     `json:"from"`
	Message *Message `json:"message"`
	Data    string   `json:"data"`
}

type Update struct {
	UpdateID      int64          `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type InlineButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboard struct {
	InlineKeyboard [][]InlineButton `json:"inline_keyboard"`
}

type apiResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
}

// --- Методы ---

func (c *Client) SendMessage(ctx context.Context, chatID, threadID int64, html string) (*Message, error) {
	var msg Message
	err := c.call(ctx, "sendMessage", map[string]any{
		"chat_id":                  chatID,
		"message_thread_id":        nonZero(threadID),
		"text":                     html,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	}, &msg)
	return &msg, err
}

// SendFile отправляет картинку как фото, остальное (PDF) — как документ.
func (c *Client) SendFile(ctx context.Context, chatID, threadID int64, file io.Reader, fileName, mime, captionHTML string, keyboard *InlineKeyboard) (*Message, error) {
	method, field := "sendDocument", "document"
	if mime == "image/jpeg" || mime == "image/png" || mime == "image/webp" {
		method, field = "sendPhoto", "photo"
	}

	fields := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"caption":    captionHTML,
		"parse_mode": "HTML",
	}
	if threadID != 0 {
		fields["message_thread_id"] = strconv.FormatInt(threadID, 10)
	}
	if keyboard != nil {
		markup, _ := json.Marshal(keyboard)
		fields["reply_markup"] = string(markup)
	}

	var msg Message
	err := c.upload(ctx, method, fields, field, fileName, file, &msg)
	return &msg, err
}

// EditCaption меняет подпись; keyboard == nil убирает кнопки.
func (c *Client) EditCaption(ctx context.Context, chatID, messageID int64, captionHTML string, keyboard *InlineKeyboard) error {
	if keyboard == nil {
		keyboard = &InlineKeyboard{InlineKeyboard: [][]InlineButton{}}
	}
	return c.call(ctx, "editMessageCaption", map[string]any{
		"chat_id":      chatID,
		"message_id":   messageID,
		"caption":      captionHTML,
		"parse_mode":   "HTML",
		"reply_markup": keyboard,
	}, nil)
}

func (c *Client) AnswerCallback(ctx context.Context, callbackID, text string, alert bool) error {
	return c.call(ctx, "answerCallbackQuery", map[string]any{
		"callback_query_id": callbackID,
		"text":              text,
		"show_alert":        alert,
	}, nil)
}

func (c *Client) GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]Update, error) {
	var updates []Update
	err := c.call(ctx, "getUpdates", map[string]any{
		"offset":          offset,
		"timeout":         timeoutSec,
		"allowed_updates": []string{"message", "callback_query"},
	}, &updates)
	return updates, err
}

func (c *Client) GetMe(ctx context.Context) (*User, error) {
	var me User
	err := c.call(ctx, "getMe", map[string]any{}, &me)
	return &me, err
}

// --- Транспорт ---

func (c *Client) call(ctx context.Context, method string, params map[string]any, out any) error {
	for k, v := range params {
		if v == nil {
			delete(params, k)
		}
	}
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+method, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, method, out)
}

func (c *Client) upload(ctx context.Context, method string, fields map[string]string, fileField, fileName string, file io.Reader, out any) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return err
		}
	}
	part, err := w.CreateFormFile(fileField, fileName)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+method, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return c.do(req, method, out)
}

func (c *Client) do(req *http.Request, method string, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		// Не логируем URL: в нём токен бота.
		return fmt.Errorf("telegram %s: request failed", method)
	}
	defer resp.Body.Close()

	var res apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("telegram %s: decode response: %w", method, err)
	}
	if !res.OK {
		return fmt.Errorf("telegram %s: %s", method, res.Description)
	}
	if out != nil {
		return json.Unmarshal(res.Result, out)
	}
	return nil
}

func nonZero(v int64) any {
	if v == 0 {
		return nil
	}
	return v
}
