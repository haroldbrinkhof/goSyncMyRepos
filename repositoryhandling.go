package main

import (
	"fmt"
	"strings"
	"sync"
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
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	for _, activeRepository := range activeGroups {
		for _, repositoryPath := range activeRepository.repositories {
			if repositoryPath != wd {
				wg.Add(1)
				go checkoutClosestPriorCommit(when, repositoryPath, wg)
			}
		}
	}
	wg.Wait()

}

func checkoutClosestPriorCommit(when time.Time, path string, wg *sync.WaitGroup) {
	defer wg.Done()
	r, err := git.PlainOpen(path)
	logAndExitOnError(err)

	ref, err := r.Head()
	logAndExitOnError(err)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	logAndExitOnError(err)

	w, err := r.Worktree()
	logAndExitOnError(err)

	var synced bool = false
	var commitToCheckout object.Commit

	cIter.ForEach(func(c *object.Commit) error {
		if when.Equal(c.Committer.When) || when.After(c.Committer.When) {
			// account for the fact that sometimes commits seem to come out of order and only take
			// most recent one to check out
			if commitToCheckout.Committer.When.Before(c.Committer.When) {
				commitToCheckout = *c
			}
			synced = true
		}
		return nil
	})

	if synced {
		err = w.Checkout(&git.CheckoutOptions{Hash: commitToCheckout.Hash})
		fmt.Printf("syncing %s done.\n", path)
	}
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
