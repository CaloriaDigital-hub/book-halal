package middleware

import (
	"io/fs"
	"net/http"
)

// noListFileSystem wraps http.FileSystem and returns 404 for directory requests,
// preventing directory listing while still serving individual files.
type noListFileSystem struct {
	fs http.FileSystem
}

func NoListFileSystem(fs http.FileSystem) http.FileSystem {
	return noListFileSystem{fs: fs}
}

func (nfs noListFileSystem) Open(name string) (http.File, error) {
	f, err := nfs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}

	if stat.IsDir() {
		f.Close()
		return nil, fs.ErrNotExist
	}

	return f, nil
}