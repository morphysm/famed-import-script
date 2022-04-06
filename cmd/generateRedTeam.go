package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

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
	Short: "Generates a red team json from the Ethereum Foundation Vulnerability Disclosures CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("generating red team...")

		// read config from args
		cfg, err := parseArgsForGenerateRedTeam(args)
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

		// map bugs to contributors and load additional data from the GitHub API
		redTeam := Contributors{}
		client := newClient("")
		redTeam.mapBugs(client, bugs)

		// marshal contributors (redTeam) to json
		teamJSON, err := json.MarshalIndent(redTeam, "", "    ")
		if err != nil {
			return err
		}

		// write json to file
		err = ioutil.WriteFile(cfg.jsonPath, teamJSON, 0644)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateRedTeamCmd)
}
