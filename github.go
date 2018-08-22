package main

import (
	"context"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var ctx context.Context

var sourceOwner = "Colelyman"
var authorName = "Cole Lyman"
var authorEmail = "cole@colelyman.com"
var sourceRepo = "colelyman-hugo"
var branch = "master"

func CommitEntry(path string, file string) error {
	client := connectGitHub()
	repo := getRef(client)
	tree, err := getTree(path, file, client, repo)
	if err != nil {
		return err
	}
	err = pushCommit(client, repo, tree)
	if err != nil {
		return err
	}
	return nil
}

func connectGitHub() *github.Client {
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.ExpandEnv("$GIT_API_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func getRef(client *github.Client) *github.Reference {
	// refURL := strings.Split(os.ExpandEnv("$REPOSITORY_URL"), "/")
	ref, _, err := client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+branch)
	if err != nil {
		panic(err)
	}

	return ref
}

// this function adds the new file to the repo
func getTree(path string, file string, client *github.Client, ref *github.Reference) (*github.Tree, error) {
	tree, _, err := client.Git.CreateTree(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, []github.TreeEntry{github.TreeEntry{Path: github.String(path), Type: github.String("blob"), Content: github.String(file), Mode: github.String(("100644"))}})

	return tree, err
}

func pushCommit(client *github.Client, ref *github.Reference, tree *github.Tree) error {
	parent, _, err := client.Repositories.GetCommit(ctx, sourceOwner, sourceRepo, *ref.Object.SHA)
	if err != nil {
		return err
	}

	parent.Commit.SHA = parent.SHA

	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &authorName, Email: &authorEmail}
	message := "Added new micropub entry."
	commit := &github.Commit{Author: author, Message: &message, Tree: tree, Parents: []github.Commit{*parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, sourceOwner, sourceRepo, commit)
	if err != nil {
		return err
	}

	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, sourceOwner, sourceRepo, ref, false)
	return err
}
