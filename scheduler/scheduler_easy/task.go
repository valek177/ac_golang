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
	Find([]FindOperator) []T
}

type Task struct {
	uuid     UUID
	status   string
	request  []byte
	response []byte
}

type Scheduler struct {
	st   storage[Task]
	proc processor

	taskQueue  chan Task
	numWorkers int
}

func NewScheduler(st storage[any], proc processor, numWorkers, queueSize int) *Scheduler {
	// TODO
	return nil
}

func (s *Scheduler) worker() {
	// TODO
}

func (s *Scheduler) AddTask(request []byte) UUID {
	// TODO
	return ""
}

func (s *Scheduler) GetTask(uuid UUID) Task {
	// TODO
	return Task{}
}

func newUUID() UUID {
	// TODO
	return ""
}
