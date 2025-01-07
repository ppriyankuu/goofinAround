package utils

// semaphore is a simple implementation for limiting concurrency
type Semaphore struct {
	limit chan struct{}
}

// This creates a new semaphore with the given limit
func NewSemaphore(maxConcurrency int) *Semaphore {
	return &Semaphore{
		limit: make(chan struct{}, maxConcurrency),
	}
}

func (s *Semaphore) Acquire() {
	s.limit <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.limit
}
