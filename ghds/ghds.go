package ghds

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"text/tabwriter"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ReleaseHistory struct {
	Repository   string    `json:"repository"`
	Releases     []Release `json:"releases"`
	ReleaseCount int       `json:"release_count"`
}

type ReleaseAsset struct {
	Name      string `json:"name"`
	Downloads int    `json:"download_count"`
}

type Release struct {
	Name           string         `json:"name"`
	Date           time.Time      `json:"date"`
	Assets         []ReleaseAsset `json:"assets"`
	TotalDownloads int            `json:"total_downloads"`
}

type DownloadStatsService interface {
	FetchReleaseHistory() (*ReleaseHistory, error)
	FormatDownloadStats(*ReleaseHistory) (string, error)
}

type GitHubDownloadStatsOptions struct {
	Release     string
	JsonOut     bool
	ApiEndpoint string
	Token       string
	PreRelease  bool
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

func includeGitHubRelease(r *github.RepositoryRelease, options *GitHubDownloadStatsOptions) bool {
	rName := r.GetName()
	tName := r.GetTagName()
	if (options.Release == "") || (options.Release == rName && rName != "") ||
		(options.Release == tName && tName != "") {
		if options.PreRelease == false && r.GetPrerelease() == true {
			return false
		}
		if len(r.Assets) > 0 {
			return true
		}
	}
	return false
}

func (ghds *GitHubDownloadStatsService) FetchReleaseHistory() (*ReleaseHistory, error) {
	ctx := context.TODO()
	opt := &github.ListOptions{
		PerPage: 200,
	}
	releaseList := []Release{}
	releaseCount := 0

	for {
		releases, resp, err := ghds.client.Repositories.ListReleases(ctx, ghds.owner, ghds.repo, opt)
		if err != nil {
			return nil, err
		}

		for _, r := range releases {
			if includeGitHubRelease(r, ghds.options) == true {
				downloadTotal := 0
				assets := []ReleaseAsset{}
				for _, a := range r.Assets {
					asset := ReleaseAsset{
						Name:      a.GetName(),
						Downloads: a.GetDownloadCount(),
					}
					downloadTotal += asset.Downloads
					assets = append(assets, asset)
				}

				release := Release{
					Name:           r.GetName(),
					Date:           r.GetCreatedAt().Time,
					Assets:         assets,
					TotalDownloads: downloadTotal,
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

	return &ReleaseHistory{
		Repository:   fmt.Sprintf("%s/%s", ghds.owner, ghds.repo),
		Releases:     releaseList,
		ReleaseCount: releaseCount,
	}, nil
}

func (ghds *GitHubDownloadStatsService) FormatDownloadStats(history *ReleaseHistory) (string, error) {
	if ghds.options.JsonOut {
		obj, err := json.Marshal(history)
		if err != nil {
			return "", err
		}

		return string(obj), nil

	} else {
		buf := new(bytes.Buffer)
		w := tabwriter.NewWriter(buf, 0, 4, 1, ' ', tabwriter.TabIndent)
		fmt.Fprintf(w, "Repository: %s/%s\n\n", ghds.owner, ghds.repo)
		if len(history.Releases) > 0 {
			for _, rel := range history.Releases {
				fmt.Fprintf(w, "Release: %v\tDate: %v\n", rel.Name, rel.Date)
				fmt.Fprintln(w, " ")
				fmt.Fprintf(w, " Asset:\tDownloads:\n")

				for _, asset := range rel.Assets {
					fmt.Fprintf(w, " - %v\t\t%v\n", asset.Name, asset.Downloads)
				}

				fmt.Fprintf(w, "\nTotal downloads:\t%v\n\n", rel.TotalDownloads)
				fmt.Fprintf(w, "------------------------------------------\n")
				w.Flush()
			}
		}

		return buf.String(), nil
	}
}

func Build(dss DownloadStatsService) (string, error) {
	history, err := dss.FetchReleaseHistory()
	if err != nil {
		return "", err
	}

	out, err := dss.FormatDownloadStats(history)
	if err != nil {
		return "", err
	}

	return out, nil
}
