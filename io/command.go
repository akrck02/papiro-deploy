package io

import (
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

func Wget(url string) error {
	return executeCommand("wget", url)
}

func Move(path string, destination string) error {
	return executeCommand("mv", path, destination)
}
