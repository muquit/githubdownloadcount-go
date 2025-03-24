package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/////////////////////////////////////////////////////////////////////
// Print githup download counts of pojects with Releases
// Developed with Claude AI 3.7 Sonnet, working under my guidance and
// instructions.
// muquit@muquit.com Mar-21-2025
/////////////////////////////////////////////////////////////////////

const (
	version        = "1.0.3"
	githubAPIURL   = "https://api.github.com/repos/"
	userAgent      = "githubdownloadcount-go"
	smallSpanStart = "<span style=\"font-size: 0.8em; color: gray;\">"
	smallSpanEnd   = "</span>"
	myUrl          = "https://github.com/muquit/githubdownloadcount-go"
)

type Asset struct {
	Name          string `json:"name"`
	DownloadCount int    `json:"download_count"`
}

type Release struct {
	Assets []Asset `json:"assets"`
}

func main() {
	// Define command line arguments
	var user string
	var project string
	var markdown bool
	var showVersion bool

	flag.StringVar(&user, "user", "", "Name of the github user")
	flag.StringVar(&project, "project", "", "Name of the github project")
	flag.BoolVar(&markdown, "markdown", false, "Output as markdown table")
	flag.BoolVar(&showVersion, "version", false, "Show version information")

	// Custom usage function to mimic Optimist's help format
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s v%s\n", os.Args[0], version)
		fmt.Fprintf(os.Stderr, "A program to display github download count for a project\n")
		fmt.Fprintf(os.Stderr, "Usage: %s options\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Where the options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Show version if requested
	if showVersion {
		fmt.Printf("%s v%s\n", os.Args[0], version)
		os.Exit(0)
	}

	// Check for required arguments
	if user == "" || project == "" {
		fmt.Fprintf(os.Stderr, "Error: Missing required arguments. Both --user and --project are required.\n")
		flag.Usage()
		os.Exit(1)
	}

	exitCode := showDownloadCounts(user, project, markdown)
	os.Exit(exitCode)
}

// by Claude
func formatNumber(n int) string {
	// Convert the number to a string
	str := strconv.Itoa(n)

	// If the number is less than 1000, no formatting needed
	if len(str) <= 3 {
		return str
	}

	// Start from the end and work backwards
	var result strings.Builder
	for i, char := range str {
		// Add comma every 3 digits from the right
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(char)
	}

	return result.String()
}

func showDownloadCounts(user, project string, markdown bool) int {
	url := githubAPIURL + user + "/" + project + "/releases"

	// Always print URL regardless of mode
	//	fmt.Printf("API URL: %s\n\n", url)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set User-Agent header
	req.Header.Set("User-Agent", userAgent)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Check if response is an empty array or error object
	if string(body) == "[]" {
		fmt.Println("No releases found for this project")
		return 1
	}

	// First try to parse as a single object (error case)
	var errorObj map[string]interface{}
	if err := json.Unmarshal(body, &errorObj); err == nil {
		// Check if this is an error message
		if message, exists := errorObj["message"]; exists {
			fmt.Printf("GitHub API Error: %v\n", message)
			return 1
		}
	}

	// Try to parse as array of releases
	var releases []Release
	if err := json.Unmarshal(body, &releases); err != nil {
		fmt.Printf("No valid releases found: %v\n", err)
		return 1
	}

	// Track total downloads
	totalDownloads := 0

	// Get current date for "Generated on" text
	currentTime := time.Now().Format("Jan-02-2006")

	// Display download counts
	if markdown {
		fmt.Println("# Download Counts")
		fmt.Println("| File | Downloads |")
		fmt.Println("| ---- | --------- |")
		for _, release := range releases {
			for _, asset := range release.Assets {
				fmt.Printf("| %s | %d |\n", asset.Name, asset.DownloadCount)
				totalDownloads += asset.DownloadCount
			}
		}
		// Add total downloads row
		if len(releases) > 0 {
			fs := formatNumber(totalDownloads)
			fmt.Printf("\n%sTotal Downloads: %s%s  \n", smallSpanStart, fs, smallSpanEnd)
			fmt.Printf("%sGenerated on: %s%s  \n", smallSpanStart, currentTime, smallSpanEnd)
			fmt.Printf("%sAPI Link: [%s](%s)%s  \n", smallSpanStart, url, url, smallSpanEnd)
			fmt.Printf("%sGenerated with: %s%s  \n", smallSpanStart, myUrl, smallSpanEnd)
		}
	} else {
		for _, release := range releases {
			for _, asset := range release.Assets {
				fmt.Printf("%s %d\n", asset.Name, asset.DownloadCount)
				totalDownloads += asset.DownloadCount
			}
		}
		if len(releases) > 0 {
			fmt.Printf("\nTotal Downloads: %d\n", totalDownloads)
		}
		fmt.Printf("\nGenerated on: %s\n", currentTime)
		fmt.Println("Generated with: https://github.com/muquit/githubdownloadcount-go\n")
	}

	// Return exit code based on download count
	if totalDownloads > 0 {
		return 0 // Success
	}

	fmt.Println("No downloads found")
	return 1 // No downloads found
}
