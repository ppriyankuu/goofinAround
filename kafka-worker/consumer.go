package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/segmentio/kafka-go"
)

const workerCount = 3

func consumeTasks() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "tasks",
		GroupID: "task-group",
	})

	defer reader.Close()

	var wg sync.WaitGroup
	tasks := make(chan string, workerCount)

	// worker pool
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range tasks {
				log.Printf("Worker %d processing: %s\n", workerID, task)
			}
		}(i)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutdown signal received")
		cancel()
	}()

	// kafka consumer
	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("Exiting Kafka consumer loop")
				break
			}
			log.Fatalf("Failed to read message: %v", err)
		}

		task := string(message.Value)
		log.Printf("Consumed: %s\n", task)
		tasks <- task
	}

	close(tasks)
	wg.Wait()
}
