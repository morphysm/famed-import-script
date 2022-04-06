package cmd

import (
	"errors"
	"log"
	"strings"
)

type bug struct {
	affectedClients string
	uID             string
	bug             string
	bugType         string
	summary         string
	links           string
	reportedDate    string
	fixedDate       string
	publishedDate   string
	severity        string
	bountyHunter    string
	bountyPoints    string
}

// mapBugs maps data to a slice of bugs.
func mapBugs(data [][]string) []bug {
	var bugs []bug
	for i := 1; i < len(data); i++ {
		line := data[i]
		bug := bug{
			affectedClients: line[0],
			uID:             line[1],
			bug:             line[2],
			bugType:         line[3],
			summary:         line[4],
			links:           line[5],
			reportedDate:    line[7],
			fixedDate:       line[8],
			publishedDate:   line[9],
			severity:        line[11],
			bountyHunter:    line[13],
			bountyPoints:    line[14],
		}

		if bug.severity == "" {
			log.Printf("Skipped bug with UID: %s due to missing severity", bug.uID)
			continue
		}
		parsedSeverity, err := parseSeverity(bug.severity)
		if err != nil {
			log.Printf("Skipped bug with UID: %s due to unknown severity: %s", bug.uID, bug.severity)
			continue
		}
		bug.severity = parsedSeverity

		if bug.bountyPoints == "" {
			log.Printf("Skipped bug with UID: %s due to missing bounty points", bug.uID)
			continue
		}
		if bug.fixedDate == "" {
			log.Printf("Skipped bug with UID: %s due to missing fixed date", bug.uID)
			continue
		}
		if bug.reportedDate == "" {
			log.Printf("Skipped bug with UID: %s due to missing reported date", bug.uID)
			continue
		}
		if bug.reportedDate == "" {
			log.Printf("Skipped bug with UID: %s due to missing published date", bug.uID)
			continue
		}

		bugs = append(bugs, bug)
	}

	return bugs
}

// parseSeverity returns a severity parsed from a string.
func parseSeverity(severity string) (string, error) {
	if strings.EqualFold("Note", severity) {
		return "info", nil
	}
	if strings.EqualFold("Low", severity) {
		return "low", nil
	}
	if strings.EqualFold("Medium", severity) {
		return "medium", nil
	}
	if strings.EqualFold("High", severity) {
		return "high", nil
	}
	if strings.EqualFold("Critical", severity) {
		return "critical", nil
	}

	return "", errors.New("invalid severity")
}
