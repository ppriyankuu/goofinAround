package scraper

import "fmt"

func PrintResults(results map[string]map[string]interface{}) {
	for url, data := range results {
		fmt.Printf("URL: %s\n", url)
		fmt.Printf("	Title: %s\n", data["title"])
		fmt.Printf("	Description: %s\n", data["description"])

		// Printing links one per line
		fmt.Println("Links:")
		if links, ok := data["links"].([]string); ok {
			for _, link := range links {
				fmt.Printf("  - %s\n", link)
			}
		} else {
			fmt.Println("  No links found or an error occurred.")
		}

		if data["error"] != nil {
			fmt.Printf("	Error: %s\n", data["error"])
		}
		fmt.Println("-----------------------------------------------------------")
	}
}
