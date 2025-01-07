package main

import (
	"fmt"
	"sync"
)

func main() {
	files := []string{"file01.txt", "file02.txt", "file03.txt", "file04.txt"}

	var totalWordCount int

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)

		go func(file string) {
			defer wg.Done()

			wordCount, err := CountWordsInFile("files/" + file)
			if err != nil {
				fmt.Printf("Error processing the file %s: %v\n", file, err)
			}

			fmt.Printf("File %s, Words: %d\n", file, wordCount)

			mutex.Lock()
			totalWordCount += wordCount
			mutex.Unlock()
		}(file)
	}

	wg.Wait()
	fmt.Printf("Total words across all file: %d\n", totalWordCount)
}
