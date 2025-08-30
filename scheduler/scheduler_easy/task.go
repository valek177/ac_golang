//go:build task_template

package main

import (
	"crypto/rand"
	"fmt"
)

type Processor interface {
	Process([]byte) ([]byte, error)
}

type (
	UUID string
	Hash string
)

type Repository interface {
	Store(Task) UUID
	GetByUUID(UUID) Task
}

type Task struct {
	// TODO
}

const (
	StatusProcessing = "processing"
	StatusQueued     = "queued"
	StatusDone       = "done"
	StatusError      = "error"
)

type Scheduler struct {
	// TODO
}

// TODO scheduler methods
func NewScheduler(repository Repository, processor Processor, numWorkers, queueSize int, generateUUID func() UUID) (*Scheduler, error) {
	// TODO
	return nil, nil
}

func (s *Scheduler) worker() {
	// TODO
}

func (s *Scheduler) AddTask(request []byte) (UUID, error) {
	// TODO
	return "", nil
}

func (s *Scheduler) GetTask(uuid UUID) Task {
	// TODO
	return Task{}
}

func (s *Scheduler) Close() {
	// TODO
}

// Для генерации UUID можно использовать готовую функцию
func generateUUID() UUID {
	// pseudo uuid
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return UUID(uuid)
}
