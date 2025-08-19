package scheduler

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"slices"
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
	hash     string
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

	mu     sync.RWMutex
	closed bool

	closeDoneCh chan struct{}
}

func NewScheduler(st storage[Task], proc processor, numWorkers, queueSize int) (*Scheduler, error) {
	if numWorkers <= 0 {
		return nil, fmt.Errorf("incorrect workers number")
	}
	if queueSize <= 0 {
		return nil, fmt.Errorf("incorrect queue size")
	}

	scheduler := &Scheduler{
		st:          st,
		proc:        proc,
		taskQueue:   make(chan Task, queueSize),
		closeDoneCh: make(chan struct{}),
	}

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(numWorkers)

		for range numWorkers {
			go func() {
				defer wg.Done()
				scheduler.worker()
			}()
		}

		wg.Wait()
		close(scheduler.closeDoneCh)
	}()

	return scheduler, nil
}

func (s *Scheduler) worker() {
	var err error

	for t := range s.taskQueue {
		t.status = StatusProcessing

		s.st.Store(t)

		t.response, err = s.proc.Process(t.request)
		if err != nil {
			t.status = StatusError
		} else {
			t.status = StatusDone
		}

		s.st.Store(t)
	}
}

func (s *Scheduler) AddTask(request []byte) (UUID, error) {
	t := Task{
		uuid:    newUUID(),
		status:  StatusQueued,
		request: request,
		hash:    generateHash(request),
	}

	query := FindOperator{
		key:      "hash",
		operator: "equals",
		value:    t.hash,
	}

	// Не добавляем таску, если уже есть таска с таким hash & bytes
	if storageTasks := s.st.Find([]FindOperator{query}); storageTasks != nil {
		for _, v := range storageTasks {
			if slices.Equal(v.request, t.request) {
				return v.uuid, nil
			}
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return "", fmt.Errorf("scheduler pool is closed")
	}

	select {
	case s.taskQueue <- t:
		return t.uuid, nil
	default:
		return "", fmt.Errorf("scheduler pool is full")
	}
}

func (s *Scheduler) GetTask(uuid UUID) Task {
	task := s.st.Get(uuid)

	if task.uuid == "" {
		return Task{}
	}

	return task
}

// Генератор UUID
func newUUID() UUID {
	return "d97976cc-35f8-44cb-91f9-fa47a85db34b"
}

func generateHash(request []byte) string {
	h := sha256.New()

	h.Write(request)

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
