package cmd

import (
	"context"
	"net/http"

	"github.com/google/go-github/v41/github"
	"github.com/morphysm/famed-github-backend/pkg/pointer"
	"golang.org/x/oauth2"
)

// gitHubClient represents a GitHub client.
type gitHubClient struct {
	client *github.Client
}

// newClient returns a new gitHubClient.
func newClient(apiToken string) *gitHubClient {
	ctx := context.Background()
	var httpClient *http.Client
	if apiToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: apiToken},
		)
		httpClient = oauth2.NewClient(ctx, ts)
	} else {
		httpClient = &http.Client{}
	}

	return &gitHubClient{client: github.NewClient(httpClient)}
}

// postIssue posts a GitHub issue to a GitHub repository.
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

// label represents a GitHub label.
type label struct {
	name  string
	color string
}

// postLabel posts a label to a GitHub repository.
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

// getUser gets a GitHub user.
func (gHC *gitHubClient) getUser(login string) (*github.User, error) {
	ctx := context.Background()
	user, _, err := gHC.client.Users.Get(ctx, login)
	if err != nil {
		return nil, err
	}

	return user, nil
}
