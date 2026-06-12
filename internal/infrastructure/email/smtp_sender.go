package email

import (
	"context"
	"fmt"
	"net/smtp"

	appusers "book_halal/internal/application/users/commands"
)


var _ appusers.EmailSender = (*SMTPSender)(nil)

type SMTPSender struct {
	host     string
	port     string
	from     string
	username string
	password string
}

func NewSMTPSender(host, port, username, password, from string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) SendOTP(ctx context.Context, toEmail string, code string) error {

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Формируем заголовки и тело письма по стандарту RFC 822
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Подтверждение регистрации Book Halal\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"Ваш код подтверждения: %s\r\n"+
		"Код действителен 5 минут.\r\n", toEmail, code))

	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	// Отправляем письмо
	return smtp.SendMail(addr, auth, s.from, []string{toEmail}, msg)
}