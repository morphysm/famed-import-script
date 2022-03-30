package cmd

import (
	"log"
)

type config struct {
	csvPath  string
	owner    string
	repo     string
	apiToken string
}

func ParseArgs(args []string) config {
	if len(args) < 4 {
		log.Fatal("missing arguments: csvPath owner repo apiToken")
	}

	cfg := config{
		csvPath:  args[0],
		owner:    args[1],
		repo:     args[2],
		apiToken: args[3],
	}

	if cfg.csvPath == "" {
		log.Fatal("missing arguments: path")
	}

	if cfg.owner == "" {
		log.Fatal("missing arguments: owner")
	}

	if cfg.repo == "" {
		log.Fatal("missing arguments: repo")
	}

	if cfg.apiToken == "" {
		log.Fatal("missing arguments: apiToken")
	}

	return cfg
}
