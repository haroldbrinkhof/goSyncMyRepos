package main

import (
	"log"
)

type repositoryGroup struct {
	name         string
	repositories []string
}

func main() {

	if !isGitDir(workingDir()) {
		log.Fatal("this command must be called from inside a git repository")
	}
	var groups []repositoryGroup = loadConfigFile()

	var activeGroups []repositoryGroup = findActiveRepositoryGroup(groups)
	if len(activeGroups) == 0 {
		log.Fatal("can not find a group where current directory(" + workingDir() + ") belongs to")
	}

	syncReposOfActiveGroups(getTimeOfCurrentCommit(), activeGroups)

}
