package producer

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func StartProducer(id int, taskChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Creating a local random generator with an unique seed
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	tasks := []string{"TaskA", "TaskB", "TaskC", "TaskD"}

	for i := 0; i < 10; i++ { // Each producer generates 10 tasks
		task := tasks[r.Intn(len(tasks))] // Using the local random generator
		fmt.Printf("Producer %d produced: %s\n", id, task)
		taskChan <- task
		time.Sleep(time.Millisecond * 200) // Simulating task generation delay
	}

	fmt.Printf("Producer %d finished producing tasks\n", id)
}
