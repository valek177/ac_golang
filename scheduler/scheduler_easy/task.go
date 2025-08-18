//go:build task_template

package main

type processor interface {
	Process([]byte) ([]byte, error)
}

type FindOperator struct {
	key, operator, value string
}

type UUID string

type storage[T any] interface {
	Store(T) UUID
	Get(UUID) T
	Find([]FindOperator) []T // select
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
func NewScheduler(st storage[Task], proc processor, numWorkers, queueSize int) (*Scheduler, error) {
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
