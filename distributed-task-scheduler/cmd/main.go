package main

import (
	"fmt"
	"sync"
	"task-scheduler/internal/config"
	"task-scheduler/internal/producer"
	"task-scheduler/internal/result"
)

func main() {
	// Loading configuration
	cfg := config.LoadConfig()

	// Creating channels
	taskChan := make(chan string, cfg.TaskBufferSize)
	resultChan := make(chan string, cfg.ResultBufferSize)

	// WaitGroup to wait for all the tasks to complete
	var wg sync.WaitGroup

	// Starting producers
	for i := 1; i <= cfg.NumProducers; i++ {
		wg.Add(1)
		go producer.StartProducer(i, taskChan, &wg)
	}

	// Starting workers
	go result.StartResultCollector(resultChan)

	// Waiting for all the producers and workers to finish
	wg.Wait()

	// Closing channels to signal completion
	close(taskChan)
	close(resultChan)

	fmt.Println("All tasks processed successfully!")
}
