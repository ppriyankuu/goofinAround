package main

import (
	"concurrent-web-scraper/scraper"
)

func main() {
	urls := []string{
		"https://golang.org",
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://www.facebook.com",
	}

	results := scraper.ScrapeURLs(urls, 3)

	scraper.PrintResults(results)
	// for url, title := range results {
	// 	fmt.Printf("%s: %s\n", url, title)
	// }
}
