package books

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"book_halal/internal/application/books/commands"
	uuid "book_halal/internal/pkg"
)

func (c *BookController) UploadBook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		c.sendError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")

	if title == "" || author == "" || priceStr == "" {
		c.sendError(w, http.StatusBadRequest, "title, author, and price are required fields")
		return
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		c.sendError(w, http.StatusBadRequest, "invalid price format, must be an integer")
		return
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		c.sendError(w, http.StatusBadRequest, "invalid or missing document file")
		return
	}
	defer file.Close()

	tmpFileName := fmt.Sprintf("upload_%s_%s", uuid.New(), header.Filename)
	tmpFilePath := filepath.Join(os.TempDir(), tmpFileName)

	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		c.sendError(w, http.StatusInternalServerError, "failed to create temp file on server")
		return
	}

	if _, err := io.Copy(tmpFile, file); err != nil {
		tmpFile.Close()
		os.Remove(tmpFilePath)
		c.sendError(w, http.StatusInternalServerError, "failed to save uploaded file")
		return
	}
	tmpFile.Close()

	cmd := commands.UploadBookCommand{
		Title:       title,
		Author:      author,
		Description: description,
		Price:       price,
		TmpPDFPath:  tmpFilePath,
		BaseURL:     c.baseURL,
	}

	result, err := c.uploadHandler.Handle(r.Context(), cmd)
	if err != nil {
		os.Remove(tmpFilePath)
		c.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.sendJSON(w, http.StatusAccepted, map[string]string{
		"message": "book upload successfully started",
		"book_id": result.BookID,
	})
}