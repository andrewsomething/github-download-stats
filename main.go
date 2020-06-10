package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andrewsomething/github-download-stats/ghds"
)

var (
	owner    = flag.String("owner", "", "The GitHub repository's owner (required)")
	repo     = flag.String("repo", "", "The GitHub repository (required)")
	release  = flag.String("release", "", "The tag name of the release; excluding will list all releases")
	jsonFlag = flag.Bool("json", false, "Output in JSON")
	endpoint = flag.String("api-endpoint", "", "API endpoint for use with GitHub Enterprise")
	token    = flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub API token")
)

func main() {
	flag.Parse()
	if *owner == "" || *repo == "" {
		fmt.Println("Must set the repo and owner...")
		flag.Usage()
		os.Exit(1)
	}

	options := &ghds.GitHubDownloadStatsOptions{
		Release:     *release,
		JsonOut:     *jsonFlag,
		ApiEndpoint: *endpoint,
		Token:       *token,
	}

	dss := ghds.NewGitHubDownloadStatsService(*owner, *repo, options)
	if err := ghds.Build(dss); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
