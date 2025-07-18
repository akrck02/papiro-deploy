package io

import (
	"io"
	"net/http"
	"os"
)

func Wget(url string, path string) error {

	out, error := os.Create(path)
	if nil != error {
		return error
	}
	defer out.Close()

	resp, error := http.Get(url)
	if nil != error {
		return error
	}

	defer resp.Body.Close()
	_, error = io.Copy(out, resp.Body)
	return error
}
