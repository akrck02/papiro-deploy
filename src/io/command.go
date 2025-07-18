package io

import (
	"os"
	"os/exec"
	"strings"
)

func ExecuteCommand(command string, args ...string) error {
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
