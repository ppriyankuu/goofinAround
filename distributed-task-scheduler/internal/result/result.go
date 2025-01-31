package result

import (
	"fmt"
	"sync"
)

func StartResultCollector(resultChan <-chan string) {
	var mu sync.Mutex
	results := make([]string, 0)

	for result := range resultChan {
		mu.Lock()
		results = append(results, result)
		mu.Unlock()
		fmt.Printf("Collected result: %s\n", result)
	}

	fmt.Println("Final results:")
	for _, res := range results {
		fmt.Println(res)
	}
}
