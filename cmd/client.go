package cmd

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

func (gHC *gitHubClient) postIssue(owner string, repo string, issue issue) error {
	issueReq := &github.IssueRequest{
		Title:     &issue.title,
		Body:      &issue.body,
		Labels:    &issue.labels,
		Assignee:  nil,
		State:     pointer.String("open"),
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

type label struct {
	name  string
	color string
}

func (gHC *gitHubClient) postLabel(owner string, repo string, label label) error {
	labelReq := &github.Label{
		Name:  &label.name,
		Color: &label.color,
	}

	ctx := context.Background()
	_, _, err := gHC.client.Issues.CreateLabel(ctx, owner, repo, labelReq)
	if err != nil {
		return err
	}

	return nil
}

func (gHC *gitHubClient) getUser(login string) (*github.User, error) {
	ctx := context.Background()
	user, _, err := gHC.client.Users.Get(ctx, login)
	if err != nil {
		return nil, err
	}

	return user, nil
}
