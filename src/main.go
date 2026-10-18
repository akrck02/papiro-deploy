package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akrck02/papiro-deploy/io"
)

type ActionInputName string

const (
	titleInput             = "INPUT_TITLE"
	descriptionInput       = "INPUT_DESCRIPTION"
	logoInput              = "INPUT_LOGO"
	accentInput            = "INPUT_ACCENT"
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
const logoFile = "resources/icons/logo.svg"

type ActionInput struct {
	title             string
	description       string
	logo              string
	accent            string
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
	Accent           string            `json:"accent,omitempty"`
	Path             map[string]string `json:"path"`
}

func main() {

	input := loadActionInput()
	showTitle(&input)

	error := moveFilesToTempDir()
	if nil != error {
		fail(error)
	}

	error = getLatestPapiro()
	if nil != error {
		fail(error)
	}

	error = movePapiroToRoot()
	if nil != error {
		fail(error)
	}

	error = copyLogo(&input)
	if nil != error {
		fail(error)
	}

	error = indexFiles(&input)
	if nil != error {
		fail(error)
	}

	error = changeConfigurationFile(&input)
	if nil != error {
		fail(error)
	}

}

func loadActionInput() ActionInput {
	return ActionInput{
		title:             os.Getenv(titleInput),
		description:       os.Getenv(descriptionInput),
		logo:              os.Getenv(logoInput),
		accent:            strings.TrimSpace(os.Getenv(accentInput)),
		path:              os.Getenv(pathInput),
		isObsidianProject: os.Getenv(isObsidianProjectInput) == "true",
		showFooter:        os.Getenv(showFooterInput) == "true",
		showBreadcrumb:    os.Getenv(showBreadcrumbInput) == "true",
		showStartPage:     os.Getenv(showStartPageInpu) == "true",
	}
}

func showTitle(input *ActionInput) {

	fmt.Println("--------------------------------------")
	fmt.Println("         🚀 Papiro deploy 🚀          ")
	fmt.Println("--------------------------------------")
	fmt.Println()
	fmt.Println("title:", input.title)
	fmt.Println("description:", input.description)
	fmt.Println("logo:", input.logo)
	fmt.Println("accent:", input.accent)
	fmt.Println("path is", input.path)
	fmt.Println("obsidian:", input.isObsidianProject)
	fmt.Println("show footer:", input.showFooter)
	fmt.Println("show breadcrumb:", input.showBreadcrumb)
	fmt.Println("show start page:", input.showStartPage)
	fmt.Println()
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

// copyLogo replaces the papiro logo with the project one, if given.
// The repository files are still in ./temp at this point.
func copyLogo(input *ActionInput) error {
	if "" == strings.TrimSpace(input.logo) {
		return nil
	}

	if ".svg" != strings.ToLower(filepath.Ext(input.logo)) {
		return fmt.Errorf("the logo must be an svg file: %s", input.logo)
	}

	content, error := os.ReadFile(filepath.Join("temp", input.logo))
	if nil != error {
		return fmt.Errorf("Failed to read the logo: %s", error.Error())
	}

	return os.WriteFile(logoFile, content, 0644)
}

func indexFiles(input *ActionInput) error {
	path := ""

	if "." != input.path {
		path = input.path
	}

	// the release brings papiro's own documents, they must not be published
	error := os.RemoveAll("./resources/wiki")
	if nil != error {
		return error
	}

	error = io.Index(fmt.Sprintf("./temp/%s", path), "./resources/wiki", input.isObsidianProject)
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
	configuration.Accent = input.accent
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

// fail stops the action with an error, so the workflow run is marked as failed
func fail(error error) {
	fmt.Printf("ERROR: %s\n", error.Error())
	os.Exit(1)
}
