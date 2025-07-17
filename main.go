package main

import (
	"fmt"
	"os"

	"github.com/akrck02/papiro-deploy/io"
)

const inputPath = "INPUT_PATH"
const inputIsObsidianProject = "INPUT_ISOBSIDIANPROJECT"
const latestPapiroReleaseUrl = "https://github.com/akrck02/papiro/archive/refs/tags/latest.tar.gz"
const latestPapiroReleaseFileName = "latest.tar.gz"

func main() {
	error := getLatestPapiro()
	if nil != error {
		fmt.Printf("ERROR: %s", error.Error())
		return
	}

	error = indexFiles()
	if nil != error {
		fmt.Printf("ERROR: %s", error.Error())
		return
	}

	error = movePapiroToRoot()
	if nil != error {
		fmt.Printf("ERROR: %s", error.Error())
		return
	}
}

func getLatestPapiro() error {

	// Download the "latest" tag from github
	error := io.Wget(latestPapiroReleaseUrl, latestPapiroReleaseFileName)
	if nil != error {
		return fmt.Errorf("Failed to download latest papiro version: %s", error.Error())
	}

	// uncompress the website
	error = io.Untar(latestPapiroReleaseFileName, ".")
	if nil != error {
		return fmt.Errorf("Failed to uncompress the latest papiro version: %s", error.Error())
	}

	return nil
}

func indexFiles() error {
	path := os.Getenv(inputPath)
	isObsidianProject := os.Getenv(inputIsObsidianProject) == "true"
	println(fmt.Sprint("path is ", path))
	println(fmt.Sprint("obsidian:", isObsidianProject))
	return nil
}

func movePapiroToRoot() error {
	// Move the files to root
	error := io.Move("papiro-latest", ".")
	if nil != error {
		return fmt.Errorf("Failed to move files to root: %s.", error.Error())
	}

	return nil
}
