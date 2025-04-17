package main

import (
	"container/heap"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	StatusQueued    = "queued"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusRetrying  = "retrying"
)

// Task represents a unit of work with dependencies, priority, and retry logic
type Task struct {
	ID           string   // Unique identifier for the task
	Data         string   // Data associated with the task
	Dependencies []string // List of task IDs this task depends on
	Status       string   // Current status of the task
	Priority     int      // Priority level (higher number means higher priority)
	RetryCount   int      // Number of retry attempts
	index        int      // Index in the priority queue (used internally by heap.Interface)
}

// TaskStatus manages the state of all tasks with thread-safe access
type TaskStatus struct {
	Tasks map[string]*Task // Map of Task ID to task object
	mu    sync.RWMutex     // Mutex for safe concurrent access
}

// Metrics tracks statistics across task execution
type Metrics struct {
	Processed int        // Total tasks successfully processed
	Failed    int        // Total tasks that failed after retries
	Retries   int        // Total retry attempts
	mu        sync.Mutex // Mutex for safe concurrent access
}

// IncrementProcessed increments the count of successfully processed tasks
func (m *Metrics) IncrementProcessed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Processed++
}

// IncrementFailed increments the count of failed tasks
func (m *Metrics) IncrementFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Failed++
}

// IncrementRetries increments the count of retry attempts
func (m *Metrics) IncrementRetries() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Retries++
}

// PriorityQueue implements heap.Interface for Tasks based on their priority
type PriorityQueue []*Task

func (pq PriorityQueue) Len() int { return len(pq) }

// higher priority tasks come first
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority // higher priority first
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// push a task into the priority queue
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	task := x.(*Task)
	task.index = n
	*pq = append(*pq, task)
}

// removes the highest priority task from the queue
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	task := old[n-1]
	old[n-1] = nil
	task.index = -1
	*pq = old[0 : n-1]
	return task
}

// processes tasks from the taskChan and sends results to resultChan
func StartWorker(id int, taskChan chan *Task, resultChan chan *Task, status *TaskStatus, metrics *Metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		// Check dependencies
		status.mu.RLock()
		dependenciesMet := true
		for _, depID := range task.Dependencies {
			if depTask, ok := status.Tasks[depID]; !ok || depTask.Status != StatusCompleted {
				dependenciesMet = false
				break
			}
		}
		status.mu.RUnlock()

		if !dependenciesMet {
			resultChan <- task // requeue
			continue
		}

		log.Printf("Worker %d processing task %s", id, task.ID)
		time.Sleep(500 * time.Millisecond) // simulate work
		success := processTask(task)

		if success {
			task.Status = StatusCompleted
			metrics.IncrementProcessed()
		} else {
			if task.RetryCount < 3 {
				task.RetryCount++
				metrics.IncrementRetries()
				log.Printf("Worker %d retrying task %s", id, task.ID)
			} else {
				task.Status = StatusFailed
				metrics.IncrementFailed()
				log.Printf("Worker %d failed task %s after retries", id, task.ID)
			}
		}

		// Update status
		status.mu.Lock()
		status.Tasks[task.ID] = task
		status.mu.Unlock()

		// Send it back to resultChan if needs retry
		if task.Status == StatusQueued || task.Status == StatusRetrying {
			resultChan <- task
		}
	}
}

func processTask(task *Task) bool {
	// Simulate some failures
	if task.ID == "task_5" || task.ID == "task_10" {
		return false
	}
	return true
}

func SubmitTask(task *Task, pq *PriorityQueue, status *TaskStatus, mu *sync.Mutex) {
	status.mu.Lock()
	status.Tasks[task.ID] = task
	status.mu.Unlock()

	mu.Lock()
	heap.Push(pq, task)
	mu.Unlock()
}

func main() {
	taskChan := make(chan *Task, 10)
	resultChan := make(chan *Task, 10)
	status := &TaskStatus{Tasks: make(map[string]*Task)}
	metrics := &Metrics{}
	var wg sync.WaitGroup

	// Priority queue
	pq := &PriorityQueue{}
	heap.Init(pq)
	var pqMu sync.Mutex

	// Launch worker pool
	workerCount := 5
	for i := range workerCount {
		wg.Add(1)
		go StartWorker(i, taskChan, resultChan, status, metrics, &wg)
	}

	// Submit initial tasks
	for i := range 20 {
		task := &Task{
			ID:       fmt.Sprintf("task_%d", i),
			Data:     fmt.Sprintf("sample_data_%d", i),
			Status:   StatusQueued,
			Priority: i % 5,
		}
		SubmitTask(task, pq, status, &pqMu)
	}

	// Scheduler
	go func() {
		for {
			pqMu.Lock()
			if pq.Len() > 0 {
				task := heap.Pop(pq).(*Task)
				pqMu.Unlock()
				taskChan <- task
			} else {
				pqMu.Unlock()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Requeue failed/retry tasks from workers
	go func() {
		for task := range resultChan {
			pqMu.Lock()
			heap.Push(pq, task)
			pqMu.Unlock()
		}
	}()

	// Wait for some time, then stop
	time.Sleep(15 * time.Second)
	close(taskChan)
	close(resultChan)
	wg.Wait()

	log.Printf("Processed: %d, Failed: %d, Retries: %d", metrics.Processed, metrics.Failed, metrics.Retries)
}
