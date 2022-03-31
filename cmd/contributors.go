package cmd

import (
	"errors"
	"log"
	"math"
	"sort"
	"time"
)

type Contributors map[string]*Contributor

type Contributor struct {
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

func (contributors Contributors) mapBug(client *gitHubClient, id string, bountyHunter string, reportedDate, publishedDate time.Time, reward float64, severity string) {
	// Get red teamer from map
	redT, ok := contributors[bountyHunter]
	if !ok {
		redT = &Contributor{
			Severities: map[string]int{},
		}
		contributors[bountyHunter] = redT

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

func (redT *Contributor) mapBug(reportedDate, publishedDate time.Time, reward float64, severity string) {
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

func (redT *Contributor) addUserIcon(client *gitHubClient) error {
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
func (contributors Contributors) updateMeanAndDeviationOfDisclosure() {
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
func (contributors Contributors) updateAverageSeverity() {
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

func (contributors Contributors) toSortedSlice() []*Contributor {
	contributorsSlice := contributors.toSlice()
	sortContributors(contributorsSlice)
	return contributorsSlice
}

// mapToSlice transforms the contributors map to a contributors slice.
func (contributors Contributors) toSlice() []*Contributor {
	contributorsSlice := make([]*Contributor, 0)
	for _, contributor := range contributors {
		contributorsSlice = append(contributorsSlice, contributor)
	}

	return contributorsSlice
}

// sortContributors sorts the contributors by descending reward sum.
func sortContributors(contributors []*Contributor) {
	sort.SliceStable(contributors, func(i, j int) bool {
		return contributors[i].RewardSum > contributors[j].RewardSum
	})
}
