package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var githubLogins = map[string]string{
	"Jonny Rhea":                 "jrhea",
	"Alexander Sadovskyi":        "AlexSSD7",
	"Martin Holst Swende":        "holiman",
	"Tintin":                     "tintinweb",
	"Antoine Toulme":             "atoulme",
	"Stefan Kobrc":               "",
	"Quan":                       "cryptosubtlety",
	"WINE Academic Workshop":     "",
	"Proto":                      "protolambda",
	"Taurus":                     "",
	"Saulius Grigaitis (+team).": "sifraitech",
	"Antonio Sanso":              "asanso",
	"Guido Vranken":              "guidovranken",
	"Jacek":                      "arnetheduck",
	"Onur Kılıç":                 "kilic",
	"Jim McDonald":               "mcdee",
	"Nishant (Prysm)":            "nisdas",
}

// generateRedTeamCmd represents the generateRedTeam command
var generateRedTeamCmd = &cobra.Command{
	Use:   "generateRedTeam",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("generateRedTeam called")

		cfg := ParseArgs(args)

		data, err := readCSV(cfg.csvPath)
		if err != nil {
			return err
		}

		bugs := mapBugs(data)
		RedTeam := Contributors{}
		client := newClient(cfg.apiToken)
		for _, bug := range bugs {
			// Parse dates and reward
			if bug.publishedDate == "" {
				log.Printf("skipping bug %s due to missing published date", bug.uID)
				continue
			}

			reportedDate, err := parseDate(bug.reportedDate)
			if err != nil {
				log.Printf("error while parsing reported date for bug with UID: %s: %v", bug.uID, err)
				continue
			}

			publishedDate, err := parseDate(bug.publishedDate)
			if err != nil {
				log.Printf("error while parsing reward date for bug with UID: %s: %v", bug.uID, err)
				continue
			}

			reward, err := strconv.ParseFloat(bug.bountyPoints, 64)
			if err != nil {
				log.Printf("error while parsing reward for bug with UID: %s: %v", bug.uID, err)
				continue
			}

			// Split bounty hunters if two are present seperated by ", "
			bountyHunters := strings.Split(bug.bountyHunter, ", ")
			for _, hunter := range bountyHunters {
				RedTeam.mapBug(client, bug.uID, hunter, reportedDate, publishedDate, reward/float64(len(bountyHunters)), bug.severity)
			}

		}

		// Calculate means
		RedTeam.updateMeanAndDeviationOfDisclosure()
		RedTeam.updateAverageSeverity()

		sortedTeam := RedTeam.toSortedSlice()

		teamJSON, err := json.MarshalIndent(sortedTeam, "", "    ")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile("redTeam.json", teamJSON, 0644)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateRedTeamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateRedTeamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateRedTeamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
