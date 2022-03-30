package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"
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
}
