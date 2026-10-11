package purchases

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

var errUnsupportedType = errors.New("unsupported file type")

// readUpload берёт файл из multipart-формы и определяет тип по содержимому, а не по имени.
// Возвращает file == nil без ошибки, если поле не передано.
func readUpload(r *http.Request, field string, allowed map[string]string) (multipart.File, string, error) {
	file, _, err := r.FormFile(field)
	if errors.Is(err, http.ErrMissingFile) {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}

	head := make([]byte, 512)
	n, _ := io.ReadFull(file, head)
	mime := http.DetectContentType(head[:n])
	if _, ok := allowed[mime]; !ok {
		file.Close()
		return nil, "", errUnsupportedType
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		file.Close()
		return nil, "", err
	}
	return file, mime, nil
}
