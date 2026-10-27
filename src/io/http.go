package io

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// downloadAttempts and retryWait give GitHub time to serve a release that
// was just created again: until then it can answer with an error page
const downloadAttempts = 5

var retryWait = 10 * time.Second

// Wget downloads url into path, trying again when the server fails
func Wget(url string, path string) error {
	var err error
	for attempt := 1; attempt <= downloadAttempts; attempt++ {
		err = download(url, path)
		if err == nil {
			return nil
		}
		fmt.Printf("download failed (attempt %d of %d): %s\n", attempt, downloadAttempts, err.Error())
		if attempt < downloadAttempts {
			time.Sleep(retryWait)
		}
	}
	return err
}

func download(url string, path string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s answered %s", url, response.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	return err
}
