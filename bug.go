package main

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

func mapBugs(data [][]string) []bug {
	var bugs []bug
	for i := 1; i < len(data); i++ {
		line := data[i]
		bugEl := bug{
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
		if bugEl.severity == "" || bugEl.bountyPoints == "" {
			continue
		}

		bugs = append(bugs, bugEl)
	}

	return bugs
}
