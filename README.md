# famed-import-script

The Famed import script imports a CSV and allows you to:
1. Generate and post GitHub issues to a repository.
2. Generate a red team json to be hosted on the Famed backend server.

## Commands

**Prerequisites**

- Go installed on your computer.
- GitHub personal access key as described [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) with access to repositories
- GitHub repository

**Arguments**

| Argument | Description                                                                                                                                                                                 | 
|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| csvPath  | path to a the CSV export                                                                                                                                                                    | 
| jsonPath | path for json output                                                                                                                                                                        |                                                                                                                                                                              |
| owner    | owner of the GitHub repository                                                                                                                                                              |
| repo     | name of the GitHub repository                                                                                                                                                               |                                                           | 
| apiToken | GitHub personal access key as described [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) with access to repositories | 

### Post Issues

This command posts issues generated from the CSV to a GitHub repository.
Additionally, it adds the famed, severity and client labels to the repository.
The commands takes some time due to a sleep interval between each post request.
The sleep interval protects against the GitHub Rate Limits.

**From Source**
Run in famed-import script folder
````
go run . postIssues <csvPath> <owner> <name> <apiToken>
````

**Compiled** 
````
postIssues <csvPath> <owner> <name> <apiToken>
````

**Note: Vulnerabilities that are missing Severity, Bounty Points, Published Date, Fixed Date or Reported Date are skipped.**

### Generate Red Team

This command generates a red team json from the CSV.
The json is intended to be hosted in the Famed Backend.
Adding it to the famed backend is currently still a manual process.

**From Source**
Run in famed-import script folder
````
go run . generateRedTeam <csvPath> <jsonPath>
````

**Compiled**
````
generateRedTeam <csvPath> <jsonPath>
````

**Note: Vulnerabilities that are missing Severity, Bounty Points, Published Date, Fixed Date or Reported Date are skipped.**