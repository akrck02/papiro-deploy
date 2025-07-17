package io

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	println(strings.Join(cmd.Args, " "))
	cmd.Dir = "."
	cmd.Stderr = os.Stderr

	error := cmd.Run()

	if error != nil {
		return error
	}

	return nil
}

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

func Move(path string, destination string) error {
	return executeCommand("mv", path, destination)
}
