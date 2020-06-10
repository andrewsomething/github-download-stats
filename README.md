# github-download-stats [![Actions Status](https://github.com/andrewsomething/github-download-stats/workflows/Test/badge.svg)](https://github.com/andrewsomething/github-download-stats/actions)

Download statistics for GitHub release assets are only available via the API. github-download-stats is a command line utility to fetch them for a specific repository.

## Installation

Pre-built binaries for Linux, macOS, and Windows are available in the [releases tab](https://github.com/andrewsomething/github-download-stats/releases).

github-download-stats can be installed from source by running:

    go get -u github.com/andrewsomething/github-download-stats

## Usage

```
Usage of ./github-download-stats:
  -api-endpoint string
    	API endpoint for use with GitHub Enterprise
  -json
    	Output in JSON
  -owner string
    	The GitHub repository's owner (required)
  -release string
    	The tag name of the release; excluding will list all releases
  -repo string
    	The GitHub repository (required)
  -token string
    	GitHub API token (default "")
  -version
    	Print version
```

## License

github-download-stats is available via the MIT license.
