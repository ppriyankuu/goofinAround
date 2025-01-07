package scraper

import (
	"concurrent-web-scraper/utils"
	"sync"
)

func ScrapeURLs(urls []string, concurrencyLimit int) map[string]map[string]interface{} {
	results := make(map[string]map[string]interface{})

	var mut sync.Mutex
	var wg sync.WaitGroup

	// Semaphore to limit concurrent goroutines
	semaphore := utils.NewSemaphore(concurrencyLimit)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			semaphore.Acquire()
			defer semaphore.Release()

			// Initialize result entry for the URL
			mut.Lock()
			results[url] = map[string]interface{}{
				"title":       "",
				"description": "",
				"error":       nil,
				"links":       nil,
			}
			mut.Unlock()

			// Fetch title and description
			title, description, err := FetchTitleAndMeta(url)
			links, linkErr := FetchLinks(url)

			mut.Lock()
			defer mut.Unlock()

			results[url]["title"] = title
			results[url]["description"] = description
			results[url]["links"] = links

			// Accumulate errors if any
			if err != nil {
				results[url]["error"] = "Error fetching title/meta: " + err.Error()
			}
			if linkErr != nil {
				if results[url]["error"] != nil {
					results[url]["error"] = results[url]["error"].(string) + "; Error fetching links: " + linkErr.Error()
				} else {
					results[url]["error"] = "Error fetching links: " + linkErr.Error()
				}
			}
		}(url)
	}

	wg.Wait()
	return results
}
