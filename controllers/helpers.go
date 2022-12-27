package controllers

import (
	"io"
	"net/http"
)

func detectContentType(f io.ReadSeeker) (string, error) {
	b := make([]byte, 512)
	if _, err := f.Read(b); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(b)

	// rewind
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return contentType, nil
}
