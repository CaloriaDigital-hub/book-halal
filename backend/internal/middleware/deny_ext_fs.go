package middleware

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// denyExtFileSystem не отдаёт файлы с указанными расширениями.
// Нужен, чтобы исходный PDF книги (static/books/<id>/book.pdf) нельзя было скачать
// в обход оплаты — RAG читает его с диска, а не по HTTP.
type denyExtFileSystem struct {
	fs   http.FileSystem
	deny []string
}

func DenyExtensions(fs http.FileSystem, exts ...string) http.FileSystem {
	return denyExtFileSystem{fs: fs, deny: exts}
}

func (d denyExtFileSystem) Open(name string) (http.File, error) {
	ext := strings.ToLower(path.Ext(name))
	for _, deny := range d.deny {
		if ext == deny {
			return nil, fs.ErrNotExist
		}
	}
	return d.fs.Open(name)
}
