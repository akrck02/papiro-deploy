package main

import (
	"fmt"
	"os"

	"github.com/akrck02/papiro-deploy/io"
)

const inputPath = "INPUT_PATH"
const inputIsObsidianProject = "INPUT_ISOBSIDIANPROJECT"
const latestPapiroReleaseUrl = "https://github.com/akrck02/papiro/releases/download/latest/papiro-latest.tar.gz"
const latestPapiroReleaseFileName = "latest.tar.gz"

type ActionInput struct {
	path              string
	isObsidianProject bool
}

func main() {

	input := loadActionInput()
	error := moveFilesToTempDir(&input)
	if nil != error {
		fmt.Printf("ERROR: %s", error.Error())
		return
	}

	error = getLatestPapiro()
	if nil != error {
		fmt.Printf("ERROR: %s", error.Error())
		return
	}

	error = indexFiles(&input)
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

func loadActionInput() ActionInput {
	return ActionInput{
		path:              os.Getenv(inputPath),
		isObsidianProject: os.Getenv(inputIsObsidianProject) == "true",
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

func indexFiles(input *ActionInput) error {
	println(fmt.Sprint("path is ", input.path))
	println(fmt.Sprint("obsidian:", input.isObsidianProject))
	return nil
}

func moveFilesToTempDir(input *ActionInput) error {
	io.Move(".", "temp")
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
