package cmd

import (
	"errors"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Contributors represents a map of Contributors.
type Contributors map[string]*Contributor

// Contributor represents a contributor to a red or blue team contributor to a software project.
type Contributor struct {
	Login            string           `json:"login"`
	AvatarURL        string           `json:"avatarUrl"`
	HTMLURL          string           `json:"htmlUrl"`
	FixCount         int              `json:"fixCount"`
	Rewards          []Reward         `json:"rewards"`
	RewardSum        float64          `json:"rewardSum"`
	Currency         string           `json:"currency"`
	RewardsLastYear  RewardsLastYear  `json:"rewardsLastYear"`
	TimeToDisclosure TimeToDisclosure `json:"timeToDisclosure"`
	Severities       map[string]int   `json:"severities"`
	MeanSeverity     float64          `json:"meanSeverity"`
}

// Reward represents a reward received for a contribution.
type Reward struct {
	Date   time.Time `json:"date"`
	Reward float64   `json:"reward"`
}

// TimeToDisclosure represents the time it took to disclose a vulnerability.
type TimeToDisclosure struct {
	Time              []float64 `json:"time"`
	Mean              float64   `json:"mean"`
	StandardDeviation float64   `json:"standardDeviation"`
}

// RewardsLastYear represents the rewards a Contributor received in the last year.
type RewardsLastYear []MonthlyReward

// MonthlyReward represents the reward a contributor received in a month.
type MonthlyReward struct {
	Month  string  `json:"month"`
	Reward float64 `json:"reward"`
}

// mapBugs maps a slice of bugs to the contributor map.
func (cs Contributors) mapBugs(client *gitHubClient, bugs []bug) {
	for _, bug := range bugs {
		// Verify and parse dates and reward
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

		// Split bounty hunters if two are present separated by ", "
		bountyHunters := strings.Split(bug.bountyHunter, ", ")
		for _, hunter := range bountyHunters {
			cs.mapBug(client, bug.uID, hunter, reportedDate, publishedDate, reward/float64(len(bountyHunters)), bug.severity)
		}
	}

	// Calculate means
	cs.updateMeanAndDeviationOfDisclosure()
	cs.updateAverageSeverity()
}

// mapBug maps a bug to the contributors map.
func (cs Contributors) mapBug(client *gitHubClient, id string, bountyHunter string, reportedDate, publishedDate time.Time, reward float64, severity string) {
	// Get red team contributor from map
	contributor, ok := cs[bountyHunter]
	if !ok {
		contributor = &Contributor{
			Severities: map[string]int{},
		}
		cs[bountyHunter] = contributor

		// Set login
		login := githubLogins[bountyHunter]
		if login != "" {
			contributor.Login = login

			// Get icons
			err := contributor.addUserIcon(client)
			if err != nil {
				log.Printf("error while retrieving user icon for bug with UID: %s: %v", id, err)
			}
		}
		if login == "" {
			log.Printf("no GitHub login found for bounty hunter %s", bountyHunter)
			contributor.Login = bountyHunter
		}

		// Set currency
		contributor.Currency = "POINTS"
	}

	contributor.mapBug(reportedDate, publishedDate, reward, severity)
}

// mapBug maps a bug to a contributor
func (c *Contributor) mapBug(reportedDate, publishedDate time.Time, reward float64, severity string) {
	// Set reward
	c.Rewards = append(c.Rewards, Reward{Date: publishedDate, Reward: reward})

	// Updated reward sum
	c.RewardSum += reward

	// Increment fix count
	c.FixCount++
	severityCount := c.Severities[severity]
	severityCount++
	c.Severities[severity] = severityCount

	// Update times to disclosure
	c.TimeToDisclosure.Time = append(c.TimeToDisclosure.Time, publishedDate.Sub(reportedDate).Minutes())
}

// addUserIcon adds a user icon to a contributor fetched from the GitHub API.
func (c *Contributor) addUserIcon(client *gitHubClient) error {
	if c.Login == "" {
		return errors.New("login is a empty string")
	}

	user, err := client.getUser(c.Login)
	if err != nil {
		return err
	}

	if user.AvatarURL != nil {
		c.AvatarURL = *user.AvatarURL
	}
	if user.HTMLURL != nil {
		c.HTMLURL = *user.HTMLURL
	}

	return nil
}

// parseDate returns a time.Time parsed from a string.
func parseDate(data string) (time.Time, error) {
	const layout = "2006-01-02"

	date, err := time.Parse(layout, data)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

// updateMeanAndDeviationOfDisclosure updates the mean and deviation of the time to disclosure of all contributors.
func (cs Contributors) updateMeanAndDeviationOfDisclosure() {
	for _, contributor := range cs {
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
func (cs Contributors) updateAverageSeverity() {
	for _, contributor := range cs {
		if contributor.FixCount == 0 {
			continue
		}

		contributor.MeanSeverity = (2*float64(contributor.Severities["low"]) +
			5.5*float64(contributor.Severities["medium"]) +
			9*float64(contributor.Severities["high"]) +
			9.5*float64(contributor.Severities["critical"])) / float64(contributor.FixCount)
	}
}

// toSortedSlice transforms the contributors map to a sorted contributors slice.
func (cs Contributors) toSortedSlice() []*Contributor {
	contributorsSlice := cs.toSlice()
	sortContributors(contributorsSlice)
	return contributorsSlice
}

// toSlice transforms the contributors map to a contributors slice.
func (cs Contributors) toSlice() []*Contributor {
	contributorsSlice := make([]*Contributor, 0)
	for _, contributor := range cs {
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
