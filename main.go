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
	getLatestPapiro()
	indexFiles()
	movePapiroToRoot()
}

func getLatestPapiro() {

	// Download the "latest" tag from github
	error := io.Wget(latestPapiroReleaseUrl, latestPapiroReleaseFileName)
	if nil != error {
		panic(fmt.Sprintf("ERROR: Failed to download latest papiro version: %s", error.Error()))
	}

	// uncompress the website
	error = io.Untar(latestPapiroReleaseFileName, ".")
	if nil != error {
		panic(fmt.Sprintf("ERROR: Failed to uncompress the latest papiro version: %s", error.Error()))
	}

}

func indexFiles() {
	path := os.Getenv(inputPath)
	isObsidianProject := os.Getenv(inputIsObsidianProject) == "true"
	println(fmt.Sprint("path is ", path))
	println(fmt.Sprint("obsidian:", isObsidianProject))
}

func movePapiroToRoot() {
	// Move the files to root
	error := io.Move("papiro-latest/*", ".")
	if nil != error {
		panic(fmt.Sprintf("ERROR: Failed to move files to root: %s.", error.Error()))
	}
}
