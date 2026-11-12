package main

import (
	"sync"
	"time"
	"strconv"
	"log"
)

type entry struct {
	value     string
	expiresAt time.Time
}

type Store struct {
	mu   sync.RWMutex
	data map[string]entry
}

func NewStore() *Store {
	return &Store{data: make(map[string]entry)}
}

func (s *Store) Set(key, value string, expiresAt time.Time) {
	s.mu.Lock()
	s.data[key] = entry{value: value, expiresAt: expiresAt}
	s.mu.Unlock()
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	e, ok := s.data[key]
	s.mu.RUnlock()

	if !ok || (!e.expiresAt.IsZero() && time.Now().After(e.expiresAt)) {
		return "", false
	}
	return e.value, true
}

func (s *Store) Delete(key []string) (int, error){
	deletedCount := 0
	for i := 0; i < len(key); i++{
		if _, exists := s.data[key[i]]; exists{
			s.mu.Lock()
			delete(s.data, key[i])
			s.mu.Unlock()
			deletedCount++
		}
	}
	return deletedCount, nil
}

func (s *Store) Exists(key []string) (int, error){
	count := 0
	for i := 0; i < len(key); i++{
		if _, exists := s.data[key[i]]; exists{
			s.mu.Lock()
			count++
			s.mu.Unlock()
		}
	}
	return count, nil
}

func (s *Store) Increment(key string) (int, bool){
	num, err := strconv.Atoi(s.data[key].value)
	if err != nil {
		log.Printf("Increment err: %v", err)
	}
	entry := s.data[key]
	num++
	entry.value = strconv.Itoa(num)
	s.mu.Lock()
	s.data[key] = entry
	s.mu.Unlock()
	return num, true
}

func (s *Store) sweep() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for key, e := range s.data {
			if !e.expiresAt.IsZero() && now.After(e.expiresAt) {
				delete(s.data, key)
			}
		}
		s.mu.Unlock()
	}
}