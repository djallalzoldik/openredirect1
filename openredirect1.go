package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// Read URLs from stdin
	scanner := bufio.NewScanner(os.Stdin)
	var urls []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			urls = append(urls, line)
		}
	}

	// Process each URL
	var results []map[string]string
	for _, urlStr := range urls {
		// extract domain from url
		urlObj, err := url.Parse(urlStr)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			continue
		}
		originalDomain := urlObj.Hostname()
		parts := strings.Split(originalDomain, ".")
		if len(parts) > 2 {
			originalDomain = strings.Join(parts[len(parts)-2:], ".")
		}

		// Perform request
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}

		// Cleanup when this function ends
		defer resp.Body.Close()

		// Check for redirection
		location := resp.Header.Get("Location")
		if location != "" {
			locationObj, err := url.Parse(location)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				continue
			}

			locationDomain := locationObj.Hostname()
			parts2 := strings.Split(locationDomain, ".")
			if len(parts2) > 2 {
				locationDomain = strings.Join(parts2[len(parts2)-2:], ".")
			}

			if locationDomain != originalDomain {
				result := map[string]string{
					"vulnerable_url": urlStr,
					"location_to": location,
				}
				results = append(results, result)
			}
		}
	}

	// Print results as JSON
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
