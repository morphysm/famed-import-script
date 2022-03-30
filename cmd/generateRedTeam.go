/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const githubHost = "https://github.com"

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

type RedTeam map[string]*RedTeamer

type RedTeamer struct {
	Login            string           `json:"login"`
	AvatarURL        *string          `json:"avatarUrl,omitempty"`
	HTMLURL          *string          `json:"htmlUrl,omitempty"`
	GravatarID       *string          `json:"gravatarId,omitempty"`
	FixCount         int              `json:"fixCount,omitempty"`
	Rewards          []Reward         `json:"rewards"`
	RewardSum        float64          `json:"rewardSum"`
	Currency         string           `json:"currency"`
	RewardsLastYear  RewardsLastYear  `json:"rewardsLastYear,omitempty"`
	TimeToDisclosure TimeToDisclosure `json:"timeToDisclosure"`
	Severities       map[string]int   `json:"severities"`
	MeanSeverity     float64          `json:"meanSeverity"`
}

type Reward struct {
	Date   time.Time `json:"date"`
	Reward float64   `json:"reward"`
}

type TimeToDisclosure struct {
	Time              []float64 `json:"time"`
	Mean              float64   `json:"mean"`
	StandardDeviation float64   `json:"standardDeviation"`
}

type RewardsLastYear []MonthlyReward

type MonthlyReward struct {
	Month  string  `json:"month"`
	Reward float64 `json:"reward"`
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
		RedTeam := RedTeam{}
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

			bountyHunters := strings.Split(bug.bountyHunter, ", ")
			for _, hunter := range bountyHunters {
				if len(bountyHunters) > 1 {
					fmt.Printf("ge")
				}
				RedTeam.mapBug(client, bug.uID, hunter, reportedDate, publishedDate, reward/float64(len(bountyHunters)), bug.severity)
			}

		}

		// Calculate means
		RedTeam.updateMeanAndDeviationOfDisclosure()
		RedTeam.updateAverageSeverity()

		sortedTeam := RedTeam.toSortedSlice()

		teamJson, err := json.MarshalIndent(sortedTeam, "", "    ")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile("redTeam.json", teamJson, 0644)
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

func (rT RedTeam) mapBug(client *gitHubClient, id string, bountyHunter string, reportedDate, publishedDate time.Time, reward float64, severity string) {
	// Get red teamer from map
	redT, ok := rT[bountyHunter]
	if !ok {
		redT = &RedTeamer{
			Severities: map[string]int{},
		}
		rT[bountyHunter] = redT

		// Set login
		login := githubLogins[bountyHunter]
		if login == "" {
			log.Printf("no login found for bounty hunter %s", bountyHunter)
		}
		redT.Login = login

		// Get icons
		err := redT.addUserIcon(client)
		if err != nil {
			log.Printf("error while retrieving user icon for bug with UID: %s: %v", id, err)
		}

		// Set currency
		redT.Currency = "POINTS"
	}

	redT.mapBug(reportedDate, publishedDate, reward, severity)
}

func (redT *RedTeamer) mapBug(reportedDate, publishedDate time.Time, reward float64, severity string) {
	// Set reward
	redT.Rewards = append(redT.Rewards, Reward{Date: publishedDate, Reward: reward})

	// Updated reward sum
	redT.RewardSum += reward

	// Increment fix count
	redT.FixCount++
	severityCount := redT.Severities[severity]
	severityCount++
	redT.Severities[severity] = severityCount

	// Update times to disclosure
	redT.TimeToDisclosure.Time = append(redT.TimeToDisclosure.Time, publishedDate.Sub(reportedDate).Minutes())
}

func (redT *RedTeamer) addUserIcon(client *gitHubClient) error {
	if redT.Login == "" {
		return errors.New("login is empty string")
	}

	user, err := client.getUser(redT.Login)
	if err != nil {
		return err
	}

	redT.AvatarURL = user.AvatarURL
	redT.GravatarID = user.GravatarID
	redT.HTMLURL = user.HTMLURL

	return nil
}

func parseDate(data string) (time.Time, error) {
	const layout = "2006-01-02"

	date, err := time.Parse(layout, data)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

// updateMeanAndDeviationOfDisclosure updates the mean and deviation of the time to disclosure of all contributors.
func (contributors RedTeam) updateMeanAndDeviationOfDisclosure() {
	for _, contributor := range contributors {
		if contributor.FixCount == 0 {
			continue
		}

		// Calculate mean
		var totalTime, sd float64
		for _, timeToDisclosure := range contributor.TimeToDisclosure.Time {
			totalTime += timeToDisclosure
		}

		contributor.TimeToDisclosure.Mean = totalTime / float64(contributor.FixCount)

		// Calculate standard deviation
		for _, timeToDisclosure := range contributor.TimeToDisclosure.Time {
			sd += math.Pow(timeToDisclosure-contributor.TimeToDisclosure.Mean, 2) //nolint:gomnd
		}

		contributor.TimeToDisclosure.StandardDeviation = math.Sqrt(sd / float64(contributor.FixCount))
	}
}

// updateAverageSeverity updates the average severity field of all contributors.
func (contributors RedTeam) updateAverageSeverity() {
	for _, contributor := range contributors {
		if contributor.FixCount == 0 {
			continue
		}

		contributor.MeanSeverity = (2*float64(contributor.Severities["low"]) +
			5.5*float64(contributor.Severities["medium"]) +
			9*float64(contributor.Severities["high"]) +
			9.5*float64(contributor.Severities["critical"])) / float64(contributor.FixCount)
	}
}

func (contributors RedTeam) toSortedSlice() []*RedTeamer {
	contributorsSlice := contributors.toSlice()
	sortContributors(contributorsSlice)
	return contributorsSlice
}

// mapToSlice transforms the contributors map to a contributors slice.
func (contributors RedTeam) toSlice() []*RedTeamer {
	contributorsSlice := make([]*RedTeamer, 0)
	for _, contributor := range contributors {
		contributorsSlice = append(contributorsSlice, contributor)
	}

	return contributorsSlice
}

// sortContributors sorts the contributors by descending reward sum.
func sortContributors(contributors []*RedTeamer) {
	sort.SliceStable(contributors, func(i, j int) bool {
		return contributors[i].RewardSum > contributors[j].RewardSum
	})
}
