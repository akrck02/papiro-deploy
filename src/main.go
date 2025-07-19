package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

type EnvironmentVariables string

const (
	githubRepository = "GITHUB_REPOSITORY"
)

const latestPapiroReleaseUrl = "https://github.com/akrck02/papiro/releases/download/latest/papiro-latest.tar.gz"
const latestPapiroReleaseFileName = "latest.tar.gz"
const configurationJsonFile = "gtdf.config.json"

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

type PapiroConfiguration struct {
	AppName          string            `json:"app_name"`
	AppDescription   string            `json:"app_description"`
	AppVersion       string            `json:"app_version"`
	CoreName         string            `json:"core_name"`
	CoreVersion      string            `json:"core_version"`
	Author           string            `json:"author"`
	GithubRepository string            `json:"github_repository"`
	WebSubpath       string            `json:"web_subpath"`
	ShowStartPage    bool              `json:"show_start_page"`
	ShowFooter       bool              `json:"show_footer"`
	ShowBreadcrumb   bool              `json:"show_breadcrumb"`
	Path             map[string]string `json:"path"`
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

	error = changeConfigurationFile(&input)
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
	error := io.CopyDirectory(".", "temp")
	if nil != error {
		return error
	}

	children, error := io.ReadDirectory("./")
	if nil != error {
		return error
	}

	for _, file := range children {
		if file.Name() != "temp" && file.Name() != ".git" {
			error = os.RemoveAll(file.Name())
			if nil != error {
				return error
			}
		}
	}

	return nil
}

func getLatestPapiro() error {

	// Download the "latest" tag from github
	error := io.Wget(latestPapiroReleaseUrl, latestPapiroReleaseFileName)
	if nil != error {
		return fmt.Errorf("Failed to download latest papiro version: %s", error.Error())
	}

	// uncompress the website
	error = io.Untar(latestPapiroReleaseFileName, "./papiro-latest")
	if nil != error {
		return fmt.Errorf("Failed to uncompress the latest papiro version: %s", error.Error())
	}

	// remove tar
	error = os.Remove(latestPapiroReleaseFileName)
	if nil != error {
		return fmt.Errorf("Failed to remove the compressed papiro version: %s", error.Error())
	}

	return nil
}

func movePapiroToRoot() error {
	error := io.CopyDirectory("papiro-latest", "./")
	if nil != error {
		return fmt.Errorf("Failed to move files to root: %s.", error.Error())
	}

	error = os.RemoveAll("papiro-latest")
	if nil != error {
		return error
	}

	return nil
}

func indexFiles(input *ActionInput) error {
	path := ""

	if "." != input.path {
		path = input.path
	}

	error := io.Index(fmt.Sprintf("./temp/%s", path), "./resources/wiki", input.isObsidianProject)
	if nil != error {
		return error
	}

	error = os.RemoveAll("./temp")
	if nil != error {
		return error
	}

	return nil
}

func changeConfigurationFile(input *ActionInput) error {

	repository := os.Getenv(githubRepository)
	repository = strings.Split(repository, "/")[1] // extract repository name
	if "" == repository {
		return fmt.Errorf("%s environment variable is not set.", githubRepository)
	}

	content, error := os.ReadFile(configurationJsonFile)
	if error != nil {
		return error
	}

	var configuration PapiroConfiguration
	error = json.Unmarshal(content, &configuration)
	if nil != error {
		return error
	}

	if "" != input.title {
		configuration.AppName = input.title
	} else {
		configuration.AppName = repository
	}

	configuration.AppDescription = input.description
	configuration.ShowFooter = input.showFooter
	configuration.ShowBreadcrumb = input.showBreadcrumb
	configuration.ShowStartPage = input.showStartPage
	configuration.WebSubpath = repository

	data, error := json.Marshal(configuration)
	if nil != error {
		return error
	}

	error = os.WriteFile(configurationJsonFile, data, 644)
	if nil != error {
		return error
	}

	return nil
}
