package main

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestIncludeRelease(t *testing.T) {
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
		actual := includeRelease(tt.input, tt.arg)
		if actual != tt.expected {
			t.Errorf("includeRelease(%v): expected %v, actual %v", tt.input, tt.expected, actual)
		}
	}
}
