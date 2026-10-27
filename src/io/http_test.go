package io

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestAnErrorPageIsNotSavedAsTheDownload(test *testing.T) {
	retryWait = 0
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	path := filepath.Join(test.TempDir(), "papiro.tar.gz")
	if err := Wget(server.URL, path); err == nil {
		test.Fatal("a 404 answer must be an error")
	}
	if _, err := os.Stat(path); err == nil {
		test.Fatal("nothing must be written when the download fails")
	}
}

func TestTheDownloadIsTriedAgainAfterAFailure(test *testing.T) {
	retryWait = 0
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		if requests == 1 {
			writer.WriteHeader(http.StatusBadGateway)
			return
		}
		writer.Write([]byte("papiro"))
	}))
	defer server.Close()

	path := filepath.Join(test.TempDir(), "papiro.tar.gz")
	if err := Wget(server.URL, path); err != nil {
		test.Fatal(err)
	}
	if content, _ := os.ReadFile(path); string(content) != "papiro" {
		test.Fatalf("got %q", content)
	}
}
