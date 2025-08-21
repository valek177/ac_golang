package main

import (
	"fmt"
	"sync"
)

// Обработчик задач
type processor interface {
	Process([]byte) ([]byte, error)
}

// Для хранилища
type FindOperator struct {
	key, operator, value string
}

type UUID string

// Интерфейс хранилища (БД)
type storage[T any] interface {
	Store(T) UUID
	Get(UUID) T
	Find([]FindOperator) []T
}

// TODO для К.
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

// TODO для К.
type Scheduler struct {
	st   storage[Task]
	proc processor

	taskQueue chan Task
}

func NewScheduler(st storage[Task], proc processor, numWorkers, queueSize int) (*Scheduler, error) {
	if numWorkers <= 0 {
		return nil, fmt.Errorf("incorrect workers number")
	}
	if queueSize <= 0 {
		return nil, fmt.Errorf("incorrect queue size")
	}

	scheduler := &Scheduler{
		st:        st,
		proc:      proc,
		taskQueue: make(chan Task, queueSize),
	}

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(numWorkers)

		for i := 0; i < numWorkers; i++ {
			go func() {
				defer wg.Done()
				scheduler.worker(i)
			}()
		}

		wg.Wait()
	}()

	return scheduler, nil
}

func (s *Scheduler) worker(id int) {
	for t := range s.taskQueue {
		t.status = StatusProcessing

		s.st.Store(t)

		response, err := s.proc.Process(t.request)
		if err != nil {
			t.status = StatusError
			t.response = nil
		} else {
			t.status = StatusDone
			t.response = response
		}

		s.st.Store(t)
	}
}

func (s *Scheduler) AddTask(request []byte) (UUID, error) {
	t := Task{
		uuid:    newUUID(),
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
	return s.st.Get(uuid)
}

func (s *Scheduler) Close() {
	close(s.taskQueue)
}

// Генератор UUID
func newUUID() UUID {
	return "d97976cc-35f8-44cb-91f9-fa47a85db34b"
}
