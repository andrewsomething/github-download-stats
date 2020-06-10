package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andrewsomething/github-download-stats/ghds"
)

var (
	version string
	commit  string

	owner       = flag.String("owner", "", "The GitHub repository's owner (required)")
	repo        = flag.String("repo", "", "The GitHub repository (required)")
	release     = flag.String("release", "", "The tag name of the release; excluding will list all releases")
	jsonFlag    = flag.Bool("json", false, "Output in JSON")
	endpoint    = flag.String("api-endpoint", "", "API endpoint for use with GitHub Enterprise")
	token       = flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub API token")
	versionFlag = flag.Bool("version", false, "Print version")
)

func main() {
	flag.Parse()

	if *versionFlag {
		if version == "" {
			version = "dev"
		}
		fmt.Printf("Version: %s\nCommit: %s\n", version, commit)
		os.Exit(0)
	}

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
	out, err := ghds.Build(dss)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(out)
}
