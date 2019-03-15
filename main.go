package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/google/go-github/github"
)

var (
	owner    = flag.String("owner", "", "The GitHub repos's owner")
	repo     = flag.String("repo", "", "The GitHub repo")
	release  = flag.String("release", "", "The tag name of the release")
	jsonFlag = flag.Bool("json", false, "Output in JSON")
)

type gitHubReleaseHistory struct {
	Repository   string          `json:"repository"`
	Releases     []gitHubRelease `json:"releases"`
	ReleaseCount int             `json:"release_count"`
}

type gitHubAsset struct {
	Name      string `json:"name"`
	Downloads int    `json:"download_count"`
}

type gitHubRelease struct {
	Name           string           `json:"name"`
	Date           github.Timestamp `json:"date"`
	Assets         []gitHubAsset    `json:"assets"`
	TotalDownloads int              `json:"total_downloads"`
}

func includeRelease(r *github.RepositoryRelease, specificRelease string) bool {
	rName := r.GetName()
	if (specificRelease == "") || (specificRelease == rName && rName != "") {
		if r.GetPrerelease() != true && len(r.Assets) > 0 {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()
	if *owner == "" || *repo == "" {
		fmt.Println("Must set the repo and owner...")
		os.Exit(1)
	}

	client := github.NewClient(nil)

	ctx := context.TODO()
	opt := &github.ListOptions{
		PerPage: 200,
	}
	releaseList := []gitHubRelease{}
	releaseCount := 0

	for {
		releases, resp, err := client.Repositories.ListReleases(ctx, *owner, *repo, opt)

		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}

		for _, r := range releases {
			if includeRelease(r, *release) == true {
				downloadTotal := 0
				assets := []gitHubAsset{}
				for _, a := range r.Assets {
					if strings.HasSuffix(a.GetName(), "sha256") != true {
						asset := gitHubAsset{
							a.GetName(),
							a.GetDownloadCount(),
						}
						downloadTotal += asset.Downloads
						assets = append(assets, asset)
					}
				}

				release := gitHubRelease{
					r.GetName(),
					r.GetCreatedAt(),
					assets,
					downloadTotal,
				}
				releaseList = append(releaseList, release)
				releaseCount++
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	releaseHistory := gitHubReleaseHistory{
		fmt.Sprintf("%s/%s", *owner, *repo),
		releaseList,
		releaseCount,
	}

	if *jsonFlag {
		obj, err := json.Marshal(releaseHistory)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(obj))
	} else {
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', tabwriter.TabIndent)
		fmt.Printf("Repository: %s/%s\n\n", *owner, *repo)
		if len(releaseList) > 0 {
			for _, rel := range releaseHistory.Releases {
				fmt.Fprintf(w, "Release: %v\tDate: %v\n", rel.Name, rel.Date)
				fmt.Fprintln(w, " ")
				fmt.Fprintf(w, " Asset:\tDownloads:\n")

				for _, asset := range rel.Assets {
					if strings.HasSuffix(asset.Name, "sha256") != true {
						fmt.Fprintf(w, " - %v\t\t%v\n", asset.Name, asset.Downloads)
					}
				}

				fmt.Fprintf(w, "\nTotal downloads:\t%v\n\n", rel.TotalDownloads)
				w.Flush()
				fmt.Println("------------------------------------------")
			}
		}
	}
}
