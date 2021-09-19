package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func isGitDir(dir string) bool {
	_, err := git.PlainOpen(dir)

	return err == nil
}

func getTimeOfCurrentCommit() time.Time {
	r, err := git.PlainOpen(workingDir())
	logAndExitOnError(err)

	ref, err := r.Head()
	logAndExitOnError(err)

	commit, err := r.CommitObject(ref.Hash())
	logAndExitOnError(err)

	return commit.Committer.When
}

func syncReposOfActiveGroups(when time.Time, activeGroups []repositoryGroup) {
	var wd string = workingDir()

	for _, activeRepository := range activeGroups {
		for _, repositoryPath := range activeRepository.repositories {
			if repositoryPath != wd {
				checkoutClosestPriorCommit(when, repositoryPath)
			}
		}
	}
}

func checkoutClosestPriorCommit(when time.Time, path string) {
	r, err := git.PlainOpen(path)
	logAndExitOnError(err)

	ref, err := r.Head()
	logAndExitOnError(err)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	logAndExitOnError(err)

	w, err := r.Worktree()
	logAndExitOnError(err)

	cIter.ForEach(func(c *object.Commit) error {
		if when.Equal(c.Committer.When) || when.After(c.Committer.When) {
			err = w.Checkout(&git.CheckoutOptions{Hash: c.Hash})
			logAndExitOnError(err)
			fmt.Println(path + " synced.")
		}
		return nil
	})
}

func findActiveRepositoryGroup(groups []repositoryGroup) []repositoryGroup {
	var activeGroups []repositoryGroup
	var wd string = workingDir()

	for _, group := range groups {
		for _, repositoryPath := range group.repositories {
			if strings.HasPrefix(wd, repositoryPath) {
				activeGroups = append(activeGroups, group)
				break
			}
		}
	}

	return activeGroups
}
