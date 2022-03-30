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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("postIssues called")

		cfg := ParseArgs(args)

		data, err := readCSV(cfg.csvPath)
		if err != nil {
			return err
		}

		bugs := mapBugs(data)

		client := newClient(cfg.apiToken)

		for _, label := range labels {
			err := client.postLabel(cfg.owner, cfg.repo, label)
			if err != nil {
				log.Printf("Error while posting label with name: %s, %v", label.name, err)
			}
			time.Sleep(4 * time.Second)
		}

		for _, bug := range bugs {
			issue := newIssue(bug)
			log.Printf("Posting bug with UID: %s", bug.uID)
			err := client.postIssue(cfg.owner, cfg.repo, issue)
			if err != nil {
				log.Printf("Error while posting bug with UID: %s, %v", bug.uID, err)
			}
			time.Sleep(4 * time.Second)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(postIssuesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postIssuesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postIssuesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
