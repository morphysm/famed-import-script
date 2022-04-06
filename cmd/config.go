package cmd

import (
	"errors"
)

// config represents a command config.
type config struct {
	csvPath  string
	jsonPath string
	owner    string
	repo     string
	apiToken string
}

// parseArgsForPostIssues returns a config parsed from command line arguments.
func parseArgsForPostIssues(args []string) (config, error) {
	if len(args) < 4 {
		return config{}, errors.New("missing arguments: csvPath owner repo apiToken")
	}

	cfg := config{
		csvPath:  args[0],
		owner:    args[1],
		repo:     args[2],
		apiToken: args[3],
	}

	if cfg.csvPath == "" {
		return config{}, errors.New("missing arguments: path")
	}

	if cfg.owner == "" {
		return config{}, errors.New("missing arguments: owner")
	}

	if cfg.repo == "" {
		return config{}, errors.New("missing arguments: repo")
	}

	if cfg.apiToken == "" {
		return config{}, errors.New("missing arguments: apiToken")
	}

	return cfg, nil
}

func parseArgsForGenerateRedTeam(args []string) (config, error) {
	if len(args) < 2 {
		return config{}, errors.New("missing arguments: csvPath jsonPath")
	}

	cfg := config{
		csvPath:  args[0],
		jsonPath: args[1],
	}

	if cfg.csvPath == "" {
		return config{}, errors.New("missing arguments: path")
	}

	return cfg, nil
}
