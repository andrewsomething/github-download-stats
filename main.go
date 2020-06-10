package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andrewsomething/github-download-stats/ghds"
)

var (
	owner    = flag.String("owner", "", "The GitHub repos's owner (required)")
	repo     = flag.String("repo", "", "The GitHub repo (required)")
	release  = flag.String("release", "", "The tag name of the release")
	jsonFlag = flag.Bool("json", false, "Output in JSON")
)

func main() {
	flag.Parse()
	if *owner == "" || *repo == "" {
		fmt.Println("Must set the repo and owner...")
		flag.Usage()
		os.Exit(1)
	}

	dss := ghds.NewGitHubDownloadStatsService(*owner, *repo, *release, *jsonFlag)
	if err := ghds.Build(dss); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
