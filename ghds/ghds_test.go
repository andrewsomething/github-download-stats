package ghds

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/github"
)

var (
	mux *http.ServeMux

	server *httptest.Server

	options *GitHubDownloadStatsOptions
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	url, _ := url.Parse(server.URL)
	endpoint := fmt.Sprintf("%s/", url.String())
	options = &GitHubDownloadStatsOptions{ApiEndpoint: endpoint}
}

func teardown() {
	server.Close()
}

func TestFetchReleaseHistory(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/foo/bar/releases", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
  {
    "tag_name": "v1.0.0",
    "name": "v1.0.0",
    "prerelease": false,
    "created_at": "2013-02-27T19:35:32Z",
    "assets": [
      {
        "name": "example.zip",
        "content_type": "application/zip",
        "download_count": 42
      },
      {
        "name": "example.tar.gz",
        "content_type": "application/zip",
        "download_count": 42
      }
    ]
  },
  {
    "tag_name": "v2.0.0",
    "name": "v2.0.0",
    "prerelease": false,
    "created_at":"2013-03-27T19:35:32Z",
    "assets": [
      {
        "name": "example.zip",
        "content_type": "application/zip",
        "download_count": 85
      }
    ]
  }
]`)
	})

	timeOne, err := time.Parse(time.RFC3339, "2013-02-27T19:35:32Z")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	timeTwo, err := time.Parse(time.RFC3339, "2013-03-27T19:35:32Z")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expected := &ReleaseHistory{
		Repository: "foo/bar",
		Releases: []Release{
			Release{
				Name: "v1.0.0",
				Date: timeOne,
				Assets: []ReleaseAsset{
					ReleaseAsset{
						Name:      "example.zip",
						Downloads: 42,
					}, ReleaseAsset{
						Name:      "example.tar.gz",
						Downloads: 42,
					},
				},
				TotalDownloads: 84,
			},
			Release{
				Name: "v2.0.0",
				Date: timeTwo,
				Assets: []ReleaseAsset{
					ReleaseAsset{
						Name:      "example.zip",
						Downloads: 85,
					},
				},
				TotalDownloads: 85,
			},
		},
		ReleaseCount: 2,
	}

	dss := NewGitHubDownloadStatsService("foo", "bar", options)
	actual, err := dss.FetchReleaseHistory()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v, expected %v", actual, expected)
	}
}

func TestFormatDownloadStats(t *testing.T) {
	timeOne, err := time.Parse(time.RFC3339, "2013-02-27T19:35:32Z")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	timeTwo, err := time.Parse(time.RFC3339, "2013-03-27T19:35:32Z")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	history := &ReleaseHistory{
		Repository: "foo/bar",
		Releases: []Release{
			Release{
				Name: "v1.0.0",
				Date: timeOne,
				Assets: []ReleaseAsset{
					ReleaseAsset{
						Name:      "example.zip",
						Downloads: 42,
					}, ReleaseAsset{
						Name:      "example.tar.gz",
						Downloads: 42,
					},
				},
				TotalDownloads: 84,
			},
			Release{
				Name: "v2.0.0",
				Date: timeTwo,
				Assets: []ReleaseAsset{
					ReleaseAsset{
						Name:      "example.zip",
						Downloads: 85,
					},
				},
				TotalDownloads: 85,
			},
		},
		ReleaseCount: 2,
	}

	dss := NewGitHubDownloadStatsService("foo", "bar", options)
	actual, err := dss.FormatDownloadStats(history)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if actual != expectedReleaseHistory {
		t.Errorf("got %v, expected %v", actual, expectedReleaseHistory)
	}
}

func TestIncludeGitHubRelease(t *testing.T) {
	var (
		emptyName = ""
		name      = "v1"
		id        = int64(1)
		isTrue    = true
	)

	var includeReleaseTests = []struct {
		input    *github.RepositoryRelease
		arg      string
		expected bool
	}{
		// No assests in release
		{&github.RepositoryRelease{}, "", false},
		// Is a pre-release, doesn't match
		{&github.RepositoryRelease{
			Name:       &emptyName,
			Prerelease: &isTrue,
			Assets:     []github.ReleaseAsset{{ID: &id}},
		}, "", false},
		// Is a pre-release, matches specific release
		{&github.RepositoryRelease{
			Name:       &emptyName,
			Prerelease: &isTrue,
			Assets:     []github.ReleaseAsset{{ID: &id}},
		}, "v1", false},
		// Release with empty name
		{&github.RepositoryRelease{
			Name:   &emptyName,
			Assets: []github.ReleaseAsset{{ID: &id}},
		}, "", true},
		// Specific release matches
		{&github.RepositoryRelease{
			Name:   &name,
			Assets: []github.ReleaseAsset{{ID: &id}},
		}, "v1", true},
		// Specific release does not match
		{&github.RepositoryRelease{
			Name:   &name,
			Assets: []github.ReleaseAsset{{ID: &id}},
		}, "v2", false},
	}

	for _, tt := range includeReleaseTests {
		actual := includeGitHubRelease(tt.input, tt.arg)
		if actual != tt.expected {
			t.Errorf("includeGitHubRelease(%v): expected %v, actual %v", tt.input, tt.expected, actual)
		}
	}
}

func TestNewGitHubDownloadStatsService(t *testing.T) {
	t.Run("passing no options should not panic", func(t *testing.T) {
		NewGitHubDownloadStatsService("digitalocean", "doctl", &GitHubDownloadStatsOptions{})
	})
}

const (
	expectedReleaseHistory = `Repository: foo/bar

Release: v1.0.0 Date: 2013-02-27 19:35:32 +0000 UTC
 
 Asset:           Downloads:
 - example.zip     42
 - example.tar.gz  42

Total downloads: 84

------------------------------------------
Release: v2.0.0 Date: 2013-03-27 19:35:32 +0000 UTC
 
 Asset:        Downloads:
 - example.zip  85

Total downloads: 85

------------------------------------------
`
)
