package main

import (
	"fmt"
	"sync"
)

type Processor interface {
	Process([]byte) ([]byte, error)
}

type UUID string

type Repository interface {
	Store(Task) UUID
	GetByUUID(UUID) Task
}

type Task struct {
	uuid     UUID
	status   string
	request  []byte
	response []byte
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

	taskQueue chan Task
}

func NewScheduler(repository Repository, processor Processor, numWorkers, queueSize int,
	generateUUID func() UUID,
) (*Scheduler, error) {
	scheduler := &Scheduler{
		repository:   repository,
		processor:    processor,
		taskQueue:    make(chan Task, queueSize),
		generateUUID: generateUUID,
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
	t := Task{
		uuid:    s.generateUUID(),
		status:  StatusQueued,
		request: request,
	}

	select {
	case s.taskQueue <- t:
		return t.uuid, nil
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
