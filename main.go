package main

import (
	"fmt"
	"os"

	"github.com/akrck02/papiro-deploy/io"
)

type ActionInputName string

const (
	titleInput             = "INPUT_TITLE"
	descriptionInput       = "INPUT_DESCRIPTION"
	logoInput              = "INPUT_LOGO"
	pathInput              = "INPUT_PATH"
	isObsidianProjectInput = "INPUT_ISOBSIDIANPROJECT"
	showFooterInput        = "INPUT_SHOWFOOTER"
	showBreadcrumbInput    = "INPUT_SHOWBREADCRUMB"
	showStartPageInpu      = "INPUT_SHOWSTARTPAGE"
)

const latestPapiroReleaseUrl = "https://github.com/akrck02/papiro/releases/download/latest/papiro-latest.tar.gz"
const latestPapiroReleaseFileName = "latest.tar.gz"

type ActionInput struct {
	title             string
	description       string
	logo              string
	path              string
	isObsidianProject bool
	showFooter        bool
	showBreadcrumb    bool
	showStartPage     bool
}

func main() {

	input := loadActionInput()
	showTitle(&input)

	error := moveFilesToTempDir()
	if nil != error {
		fmt.Printf("ERROR: %s", error.Error())
		return
	}

	error = getLatestPapiro()
	if nil != error {
		fmt.Printf("ERROR: %s", error.Error())
		return
	}

	error = movePapiroToRoot()
	if nil != error {
		fmt.Printf("ERROR: %s", error.Error())
		return
	}

	error = indexFiles(&input)
	if nil != error {
		fmt.Printf("ERROR: %s", error.Error())
		return
	}

}

func loadActionInput() ActionInput {
	return ActionInput{
		title:             os.Getenv(titleInput),
		description:       os.Getenv(descriptionInput),
		logo:              os.Getenv(logoInput),
		path:              os.Getenv(pathInput),
		isObsidianProject: os.Getenv(isObsidianProjectInput) == "true",
		showFooter:        os.Getenv(showFooterInput) == "true",
		showBreadcrumb:    os.Getenv(showBreadcrumbInput) == "true",
		showStartPage:     os.Getenv(showBreadcrumbInput) == "true",
	}
}

func showTitle(input *ActionInput) {

	fmt.Sprintln("--------------------------------------")
	fmt.Sprintln("         🚀 Papiro deploy 🚀          ")
	fmt.Sprintln("--------------------------------------")
	fmt.Sprintln()
	fmt.Sprintln("title:", input.title)
	fmt.Sprintln("description:", input.description)
	fmt.Sprintln("logo:", input.logo)
	fmt.Sprintln("path is", input.path)
	fmt.Sprintln("obsidian:", input.isObsidianProject)
	fmt.Sprintln("show footer:", input.showFooter)
	fmt.Sprintln("show breadcrumb:", input.showBreadcrumb)
	fmt.Sprintln("show start page:", input.showStartPage)
	fmt.Sprintln()
}

func moveFilesToTempDir() error {
	return io.Move(".", "temp")
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

func movePapiroToRoot() error {
	error := io.Move("papiro-latest", ".")
	if nil != error {
		return fmt.Errorf("Failed to move files to root: %s.", error.Error())
	}

	return nil
}

func indexFiles(input *ActionInput) error {
	return io.Index(fmt.Sprintf("./temp/%s", input.path), "./resources/wiki", input.isObsidianProject)
}
