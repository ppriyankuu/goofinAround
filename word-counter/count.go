package main

import (
	"bufio"
	"os"
)

// reads a file and counts the number of words
func CountWordsInFile(filepath string) (int, error) {
	// open the file
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	wordCount := 0
	var scanner *bufio.Scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) // splits the text into words

	for scanner.Scan() {
		wordCount++
	}

	return wordCount, nil
}
