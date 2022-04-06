package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
)

var labels = []label{
	{name: "famed", color: "8800ff"},
	{name: "info", color: "8800ff"},
	{name: "low", color: "8800ff"},
	{name: "medium", color: "8800ff"},
	{name: "high", color: "8800ff"},
	{name: "critical", color: "8800ff"},
	{name: "Nimbus", color: "8800ff"},
	{name: "Teku", color: "8800ff"},
	{name: "Prysm", color: "8800ff"},
	{name: "Lodestar", color: "8800ff"},
	{name: "Lighthouse", color: "8800ff"},
}

// postIssuesCmd represents the postIssues command
var postIssuesCmd = &cobra.Command{
	Use:   "postIssues",
	Short: "Posts issues generated from a disclosure csv to a GitHub repository.",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("posting issues...")

		// read config from args
		cfg, err := parseArgsForPostIssues(args)
		if err != nil {
			return err
		}

		// read csv in path
		data, err := readCSV(cfg.csvPath)
		if err != nil {
			return err
		}

		// map the csv data to the internal bug datastructure
		bugs := mapBugs(data)

		// setup new GitHub client
		client := newClient(cfg.apiToken)

		// post Famed, severity and ethereum client labels to the GitHub repo
		for _, label := range labels {
			log.Printf("Posting label: %s", label)
			err := client.postLabel(cfg.owner, cfg.repo, label)
			if err != nil {
				log.Printf("Error while posting label with name: %s, %v", label.name, err)
			}

			// sleep to avoid rate limit
			time.Sleep(4 * time.Second)
		}

		// transform bugs to issues and post them to the GitHub repo
		for _, bug := range bugs {
			issue := newIssue(bug)
			log.Printf("Posting bug with UID: %s", bug.uID)
			err := client.postIssue(cfg.owner, cfg.repo, issue)
			if err != nil {
				log.Printf("Error while posting bug with UID: %s, %v", bug.uID, err)
			}

			// sleep to avoid rate limit
			time.Sleep(4 * time.Second)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(postIssuesCmd)
}
