package main

import (
	"fmt"
	"os"
	"os/user"
	"log"
	"bufio"
	"strings"
	"path/filepath"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type repositoryGroup struct {
	name string
	repositories []string
}

func main(){


	fmt.Println(workingDir())
	fmt.Println(isGitDir(workingDir()))
	groups := loadConfigFile()

	currentGroups := findActiveRepositoryGroup(groups)
	if len(currentGroups) == 0 {
		log.Fatal("can not find a group where current directory(" + workingDir() + ") belongs to")
	}
	fmt.Println("working directory belongs to group " + currentGroups[0].name)
	fmt.Println("working directory belongs to group " + currentGroups[1].name)
	fmt.Printf("working directory belongs to %d groups\n", len(currentGroups))

	r, err := git.PlainOpen(workingDir())
	fmt.Println(err)
	ref, err := r.Head()
	fmt.Println(err)
	//CheckIfError(err)

	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	fmt.Println(err)
	//CheckIfError(err)
	fmt.Println(commit)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	fmt.Println(err)
	cCount := 0
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c.Committer)
		fmt.Println(c.Author)
		fmt.Println(c.Hash)
		cCount++

		return nil
	})
	fmt.Println(err)

}

func findActiveRepositoryGroup(groups []repositoryGroup) []repositoryGroup{
	var activeGroups []repositoryGroup

	wd := workingDir()
	for _, group := range groups {
		for _, repositoryPath := range group.repositories {
			if strings.HasPrefix(wd, repositoryPath) {
				activeGroups = append(activeGroups, group)
				break;
			}
		}
	}

	return activeGroups
}

func workingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("could not determine working directory.")
	}
	return wd
}


func isGitDir(dir string) bool {
	_, err := git.PlainOpen(dir)

	return err == nil
}

func loadConfigFile() []repositoryGroup {
	var groups []repositoryGroup;

	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(currentUser.Username + "=" + currentUser.HomeDir)

	configFile := filepath.Join(currentUser.HomeDir, "/goSyncRepositories.config")
	fmt.Println(configFile)
	file, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err){
			fmt.Println("please put your configuration in " + configFile)
		}
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	index := -1
	for scanner.Scan() {
		configLine := strings.Trim(scanner.Text()," \t")
		groupName := strings.Trim(configLine,"[]")
		if groupName != configLine{
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
			configLine = strings.Replace(configLine,"~",currentUser.HomeDir,1)
		}
		if isGitDir(configLine) {
		groups[index].repositories = append(groups[index].repositories, configLine)
		} else {
			log.Printf("not adding %s to group %s because it's not a git repository directory.", configLine, groups[index].name)
		}
		fmt.Println(groups)
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	return groups
}


