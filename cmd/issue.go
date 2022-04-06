package cmd

import (
	"strings"
)

type issue struct {
	title  string
	body   string
	labels []string
}

// newIssue returns a new issue from a bug.
func newIssue(bug bug) issue {
	title := "Famed Retroactive Rewards: " + bug.bug

	body := "**UID:** " + bug.uID + "\n\n" +
		"**Severity:** " + bug.severity + "\n\n" +
		"**Type:** " + bug.bugType + "\n\n" +
		"**Affected Clients:** " + bug.affectedClients + "\n\n" +
		"**Summary:** " + bug.summary + "\n\n" +
		"**Links:** " + bug.links + "\n\n" +
		"**Reported:** " + bug.reportedDate + "\n\n" +
		"**Fixed:** " + bug.fixedDate + "\n\n" +
		"**Published:** " + bug.publishedDate + "\n\n" +
		"**Bounty Hunter:** " + bug.bountyHunter + "\n\n" +
		"**Bounty Points:** " + bug.bountyPoints

	labels := []string{"famed", bug.severity}
	labels = append(labels, parseClients(bug.affectedClients)...)

	return issue{title: title, labels: labels, body: body}
}

// parseClients returns a slice of ethereum clients parsed from a string
func parseClients(client string) []string {
	if strings.EqualFold("Teku", client) {
		return []string{"teku"}
	}
	if strings.EqualFold("Prysm", client) {
		return []string{"prysm"}
	}
	if strings.EqualFold("Lighthouse", client) {
		return []string{"lighthouse"}
	}
	if strings.EqualFold("Lodestar", client) {
		return []string{"lodestar"}
	}
	if strings.EqualFold("All clients", client) {
		return []string{"teku", "prysm", "lighthouse", "lodestar"}
	}

	return []string{}
}
