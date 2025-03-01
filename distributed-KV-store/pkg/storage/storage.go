package storage

import "sync"

// In-memory storage layer

type Storage struct {
	MU sync.RWMutex
	DB map[string]string // used a map to store the key-value pairs
}

func NewStorage() *Storage {
	return &Storage{
		DB: make(map[string]string),
	}
}

func (s *Storage) Get(key string) (string, bool) {
	s.MU.RLock()
	defer s.MU.RUnlock()
	val, ok := s.DB[key]
	return val, ok
}

func (s *Storage) Set(key, value string) {
	s.MU.Lock()
	defer s.MU.Unlock()
	s.DB[key] = value
}

func (s *Storage) Delete(key string) {
	s.MU.Lock()
	defer s.MU.Unlock()
	delete(s.DB, key)
}
