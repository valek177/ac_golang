package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"slices"
	"sync"
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
	GetByHash(hash Hash) []Task
}

type Task struct {
	uuid     UUID
	status   string
	request  []byte
	response []byte
	hash     Hash
}

const (
	StatusProcessing = "processing"
	StatusQueued     = "queued"
	StatusDone       = "done"
	StatusError      = "error"
)

type Scheduler struct {
	repository Repository
	processor  Processor

	wg sync.WaitGroup

	generateUUID func() UUID
	generateHash func(request []byte) Hash

	taskQueue chan Task
}

func NewScheduler(repository Repository, processor Processor, numWorkers, queueSize int,
	generateUUID func() UUID, generateHash func(request []byte) Hash,
) (*Scheduler, error) {
	scheduler := &Scheduler{
		repository:   repository,
		processor:    processor,
		taskQueue:    make(chan Task, queueSize),
		generateUUID: generateUUID,
		generateHash: generateHash,
	}

	scheduler.wg.Add(numWorkers)

	for range numWorkers {
		go func() {
			defer scheduler.wg.Done()
			scheduler.worker()
		}()
	}

	return scheduler, nil
}

func (s *Scheduler) worker() {
	for task := range s.taskQueue {
		task.status = StatusProcessing

		s.repository.Store(task)

		response, err := s.processor.Process(task.request)
		if err != nil {
			task.status = StatusError
			task.response = nil
		} else {
			task.status = StatusDone
			task.response = response
		}

		s.repository.Store(task)
	}
}

func (s *Scheduler) AddTask(request []byte) (UUID, error) {
	task := Task{
		uuid:    s.generateUUID(),
		status:  StatusQueued,
		request: request,
		hash:    s.generateHash(request),
	}

	// Не добавляем таску, если уже есть таска с таким hash & bytes
	storageTasks := s.repository.GetByHash(task.hash)
	if len(storageTasks) > 0 {
		for _, v := range storageTasks {
			if slices.Equal(v.request, task.request) {
				return v.uuid, nil
			}
		}
	}

	select {
	case s.taskQueue <- task:
		return task.uuid, nil
	default:
		return "", fmt.Errorf("scheduler pool is full")
	}
}

func (s *Scheduler) GetTask(uuid UUID) Task {
	return s.repository.GetByUUID(uuid)
}

func (s *Scheduler) Close() {
	close(s.taskQueue)
	s.wg.Wait()
}

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

func generateHash(request []byte) Hash {
	h := sha256.New()

	h.Write(request)

	return Hash(base64.URLEncoding.EncodeToString(h.Sum(nil)))
}
