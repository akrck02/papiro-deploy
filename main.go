package main

import (
	"fmt"
	"os"
)

func main() {
	path := os.Getenv("INPUT_PATH")
	isObsidian := os.Getenv("INPUT_ISOBSIDIANPROJECT") == "true"
	println(fmt.Sprint("path is ", path))
	println(fmt.Sprint("obsidian:", isObsidian))
}
