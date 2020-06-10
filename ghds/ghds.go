package ghds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitHubReleaseHistory struct {
	Repository   string          `json:"repository"`
	Releases     []GitHubRelease `json:"releases"`
	ReleaseCount int             `json:"release_count"`
}

type GitHubAsset struct {
	Name      string `json:"name"`
	Downloads int    `json:"download_count"`
}

type GitHubRelease struct {
	Name           string           `json:"name"`
	Date           github.Timestamp `json:"date"`
	Assets         []GitHubAsset    `json:"assets"`
	TotalDownloads int              `json:"total_downloads"`
}

type GitHubDownloadStatsOptions struct {
	Release     string
	JsonOut     bool
	ApiEndpoint string
	Token       string
}

type GitHubDownloadStatsService struct {
	owner   string
	repo    string
	client  *github.Client
	options *GitHubDownloadStatsOptions
}

func NewGitHubDownloadStatsService(owner string, repo string, options *GitHubDownloadStatsOptions) *GitHubDownloadStatsService {
	client := github.NewClient(nil)

	if options.Token != "" {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: options.Token})
		client = github.NewClient(oauth2.NewClient(context.Background(), tokenSource))
	}

	if options.ApiEndpoint != "" {
		baseURL, err := url.Parse(options.ApiEndpoint)
		if err != nil {
			panic("invalid base URL: " + err.Error())
		}
		client.BaseURL = baseURL
	}

	return &GitHubDownloadStatsService{
		owner:   owner,
		repo:    repo,
		client:  client,
		options: options,
	}
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

func Build(ghds *GitHubDownloadStatsService) error {
	ctx := context.TODO()
	opt := &github.ListOptions{
		PerPage: 200,
	}
	releaseList := []GitHubRelease{}
	releaseCount := 0

	for {
		releases, resp, err := ghds.client.Repositories.ListReleases(ctx, ghds.owner, ghds.repo, opt)

		if err != nil {
			return err
		}

		for _, r := range releases {
			if includeRelease(r, ghds.options.Release) == true {
				downloadTotal := 0
				assets := []GitHubAsset{}
				for _, a := range r.Assets {
					if strings.HasSuffix(a.GetName(), "sha256") != true {
						asset := GitHubAsset{
							a.GetName(),
							a.GetDownloadCount(),
						}
						downloadTotal += asset.Downloads
						assets = append(assets, asset)
					}
				}

				release := GitHubRelease{
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

	releaseHistory := GitHubReleaseHistory{
		fmt.Sprintf("%s/%s", ghds.owner, ghds.repo),
		releaseList,
		releaseCount,
	}

	if ghds.options.JsonOut {
		obj, err := json.Marshal(releaseHistory)
		if err != nil {
			return err
		}
		fmt.Println(string(obj))
	} else {
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', tabwriter.TabIndent)
		fmt.Printf("Repository: %s/%s\n\n", ghds.owner, ghds.repo)
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

	return nil
}
