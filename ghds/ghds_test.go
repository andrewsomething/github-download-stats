package ghds

import (
	"testing"

	"github.com/google/go-github/github"
)

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

	t.Run("passing a valid URL should not panic", func(t *testing.T) {
		options := &GitHubDownloadStatsOptions{ApiEndpoint: "https://example.com/api/v3/"}
		NewGitHubDownloadStatsService("digitalocean", "doctl", options)
	})
}
