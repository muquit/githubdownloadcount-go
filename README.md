# GitHub Download Counter

A simple tool to display download counts for GitHub release assets. This tool allows you to easily track how many times your GitHub release assets have been downloaded and can generate output in both plain text and markdown table format.

This tool is a Go port of my original Ruby script `github-download-count` 
with enhancements for improved error handling and markdown output.

## Features

- Display download counts for all assets across all releases
- Generate markdown tables for easy inclusion in README files
- Exit with appropriate status codes (0 if downloads exist, 1 if none)
- Handles projects with no releases gracefully

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/githubdownloadcount-go.git

# Build the binary
cd githubdownloadcount-go
go build -o githubdownloadcount
```
To create binaries for various platforms:

* Look at `platforms.txt` and uncomment platforms of your choice, then type:

```bash
./go-xbuild.sh
```

## Usage

```bash
./githubdownloadcount-go --user=username --project=projectname [options]
```

### Options

```
  -markdown
        Output as markdown table
  -project string
        Name of the github project (required)
  -user string
        Name of the github user (required)
```

### Examples

Basic usage:

```bash
./githubdownloadcount --user=muquit --project=mailsend-go
```

Output as markdown table:

```bash
./githubdownloadcount --user=muquit --project=mailsend-go --markdown
```

Sample output:

```
| File | Downloads |
| ---- | --------- |
| mailsend-go_1.0.1_linux_64-bit.deb | 1234 |
| mailsend-go_1.0.1_macOS_64-bit.tar.gz | 5678 |
| mailsend-go_1.0.1_windows_64-bit.zip | 9012 |
```

## Exit Codes

- `0`: Success (downloads found)
- `1`: No downloads found or error occurred

This allows for easy integration with scripts:

```bash
./githubdownloadcount --user=muquit --project=mailsend-go
if [ $? -eq 0 ]; then
    echo "Downloads found!"
else
    echo "No downloads found or an error occurred."
fi
```

## Why This Tool?

GitHub's web interface doesn't make it easy to see total download statistics across releases. This tool provides a quick way to check how popular your releases are and can be easily integrated into CI/CD pipelines or documentation generation workflows.

## Authors

Developed with Claude AI 3.7 Sonnet, working under my guidance and instructions.

## License

MIT License - See LICENSE.txt file for details.


