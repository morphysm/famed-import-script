# famed-import-script

This script is made to use a CSV export of https://docs.google.com/spreadsheets/d/1IsDVNWpmvbHoXi5fheNtWATOCvInjINngayybEmk0MM/edit#gid=0.

The Famed import script imports the Ethereum Foundation Vulnerability Disclosures as a CSV and allows you to:
1. Generate and post GitHub issues to a repository.
2. Generate a red team json to be hosted on the Famed backend server.

## Commands

**Arguments**

| Argument | Description                                                                                                                                                                                 | 
|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| csvPath  | path to a the "Ethereum Foundation Vulnerability Disclosures" CSV export                                                                                                                    | 
| jsonPath | path for json output                                                                                                                                                                        |                                                                                                                                                                              |
| owner    | owner of the GitHub repository                                                                                                                                                              |
| repo     | name of the GitHub repository                                                                                                                                                               |                                                           | 
| apiToken | GitHub personal access key as described [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) with access to repositories | 

### Post Issues

This command posts issues generated from the "Ethereum Foundation Vulnerability Disclosures" CSV to a GitHub repository.

````
postIssues <csvPath> <owner> <name> <apiToken>
````

**Note: Vulnerabilities that are missing Severity, Bounty Points, Published Date, Fixed Date or Reported Date are skipped.**

### Generate Red Team

This command generates a red team json containing information about the red team contributors.
The json is intended to be hosted in the Famed Backend.
Adding it to the famed backend is currently still a manual process.

````
generateRedTeam <csvPath> <jsonPath>
````

**Note: Vulnerabilities that are missing Severity, Bounty Points, Published Date, Fixed Date or Reported Date are skipped.**