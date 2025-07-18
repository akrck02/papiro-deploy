package io

func Index(path string, destination string, isObsidianProject bool) error {

	if isObsidianProject {
		return ExecuteCommand("./bin/papiro-indexer", "-path", path, "-destination", destination, "-obsidian")
	}

	return ExecuteCommand("./bin/papiro-indexer", "-path", path, "-destination", destination)
}
