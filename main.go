package main

import (
	"encoding/csv"
	"log"
	"os"
)

func main() {
	cfg := ReadConfig()

	f, err := os.Open(cfg.csvPath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	bugs := mapBugs(data)

	client := newClient(cfg.apiToken)

	for _, bug := range bugs {
		err := client.postComment(cfg.owner, cfg.repo, bug.bug, bug.summary, []string{"famed", bug.severity})
		if err != nil {
			log.Print(err)
		}
	}
}
