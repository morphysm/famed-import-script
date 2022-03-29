package main

import (
	"context"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
	"golang.org/x/oauth2"
)

type gitHubClient struct {
	client *github.Client
}

func newClient(apiToken string) *gitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &gitHubClient{client: github.NewClient(tc)}
}

func (gHC *gitHubClient) postComment(owner string, repo string, title string, body string, labels []string) error {
	issueReq := &github.IssueRequest{
		Title:     pointer.String(title),
		Body:      pointer.String(body),
		Labels:    &labels,
		Assignee:  nil,
		State:     nil,
		Milestone: nil,
		Assignees: nil,
	}

	ctx := context.Background()
	_, _, err := gHC.client.Issues.Create(ctx, owner, repo, issueReq)
	if err != nil {
		return err
	}

	return nil
}
