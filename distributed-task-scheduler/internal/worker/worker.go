package worker

import (
	"fmt"
	"sync"
	"time"
)

func StartWorker(id int, taskChan <-chan string, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		fmt.Printf("Worker %d processing: %s\n", id, task)
		time.Sleep(time.Millisecond * 500) // Simulating processing delay
		result := fmt.Sprintf("Processed %s by Worker %d\n", task, id)
		resultChan <- result
	}

	fmt.Printf("Worker %d finished processing tasks.\n", id)
}
