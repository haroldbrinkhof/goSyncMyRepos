package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func workingDir() string {
	wd, err := os.Getwd()
	logMessageAndExitOnError(err, "could not determine working directory.")
	return wd
}

func homeDir() string {
	currentUser, err := user.Current()
	logAndExitOnError(err)

	return currentUser.HomeDir
}

func loadConfigFile() []repositoryGroup {
	var file *os.File = openConfigFile()
	defer closeConfigFile(file)
	return parseConfig(file)
}

func determineConfigFileLocation() string {
	return filepath.Join(homeDir(), "/goSyncRepositories.config")
}

func openConfigFile() *os.File {
	var configFilePath string = determineConfigFileLocation()
	file, err := os.Open(configFilePath)
	logAndExitOnErrorWithAdditionalAction(err, printMissingConfigFileMsg)
	return file
}

func printMissingConfigFileMsg(err error) {
	if os.IsNotExist(err) {
		fmt.Println("please put your configuration in " + determineConfigFileLocation())
	}
}

func closeConfigFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

func parseConfig(file *os.File) []repositoryGroup {
	var groups []repositoryGroup
	var scanner *bufio.Scanner = bufio.NewScanner(file)

	var index int = -1
	for scanner.Scan() {
		configLine := strings.Trim(scanner.Text(), " \t")
		groupName := strings.Trim(configLine, "[]")
		if groupName != configLine {
			currentGroup := repositoryGroup{}
			currentGroup.name = groupName
			groups = append(groups, currentGroup)
			index++
			continue
		}
		if index == -1 {
			log.Fatal("please add a section on top of your repositories in the config file to group them")
		}
		if strings.HasPrefix(configLine, "~/") {
			configLine = strings.Replace(configLine, "~", homeDir(), 1)
		}
		if isGitDir(configLine) {
			groups[index].repositories = append(groups[index].repositories, configLine)
		} else {
			log.Printf("not adding %s to group %s because it's not a git repository directory.", configLine, groups[index].name)
		}
	}

	return groups
}
