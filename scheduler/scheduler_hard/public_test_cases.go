package main

import (
	"fmt"
	"slices"
	"time"
)

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Check NewScheduler is OK with valid parameters",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 2, 1)
			if err != nil {
				return false
			}

			if scheduler.st == nil {
				return false
			}

			if scheduler.proc == nil {
				return false
			}

			if scheduler.taskQueue == nil {
				return false
			}

			return true
		},
	},
	{
		name: "Check NewScheduler is FAILED with incorrect numWorkers (0)",
		check: func() bool {
			_, err := NewScheduler(makeStorage(), makeProcessor(), 0, 1)
			if err == nil {
				return false
			}

			return true
		},
	},
	{
		name: "Check NewScheduler is FAILED with incorrect numWorkers (<0)",
		check: func() bool {
			_, err := NewScheduler(makeStorage(), makeProcessor(), -1, 1)
			if err == nil {
				return false
			}

			return true
		},
	},
	{
		name: "Check NewScheduler is FAILED with incorrect queueSize (0)",
		check: func() bool {
			_, err := NewScheduler(makeStorage(), makeProcessor(), 2, 0)
			if err == nil {
				return false
			}

			return true
		},
	},
	{
		name: "Check NewScheduler is FAILED with incorrect queueSize (<0)",
		check: func() bool {
			_, err := NewScheduler(makeStorage(), makeProcessor(), 2, -1)
			if err == nil {
				return false
			}

			return true
		},
	},
	{
		name: "Check GetTask returns not empty task",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 1, 1)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{1})
			if uuid == "" {
				return false
			}

			return true
		},
	},
	{
		name: "Check AddTask is OK",
		check: func() bool {
			// TODO fix
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 1, 1)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{1})
			if uuid == "" {
				return false
			}

			return !scheduler.isEmptyTaskQueue()
		},
	},
	{
		name: "Check AddTask is OK (different tasks, without exceeding the queue size)",
		check: func() bool {
			// TODO fix
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 1, 3)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{1})
			if uuid == "" {
				return false
			}

			uuid, err = scheduler.AddTask([]byte{2})
			if uuid == "" {
				return false
			}

			uuid, err = scheduler.AddTask([]byte{3})
			if uuid == "" {
				return false
			}

			return !scheduler.isEmptyTaskQueue()
		},
	},
	{
		name: "Check AddTask is OK (task with hash & bytes already exists)",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 1, 1)
			if err != nil {
				return false
			}

			checkUuid := newUUID()

			scheduler.st.Store(Task{
				uuid:    checkUuid,
				request: []byte{1},
				hash:    generateHash([]byte{1}),
			})

			// Добавляем ту же таску, должны вернуть тот же UUID, что у таски в storage
			uuid, err := scheduler.AddTask([]byte{1})
			if uuid != checkUuid || err != nil {
				return false
			}

			return true
		},
	},
	{
		name: "Check AddTask is not OK (pool is full)",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 1, 1)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{1})
			if uuid == "" {
				return false
			}

			uuid, err = scheduler.AddTask([]byte{2})
			if err == nil || uuid != "" {
				return false
			}

			return true
		},
	},
	{
		name: "Check Close is OK",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 1, 1)
			if err != nil {
				return false
			}

			scheduler.Close()

			return scheduler.isClosed()
		},
	},
	{
		name: "Check worker is OK (task is done)",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 1, 1)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{1})
			if err != nil {
				return false
			}

			// чтобы успел выполниться worker в горутине
			time.Sleep(time.Second)

			task := scheduler.GetTask(uuid)
			if task.uuid == "" {
				return false
			}

			if !slices.Equal([]byte{1}, task.response) {
				return false
			}

			if task.status != StatusDone {
				return false
			}

			return true
		},
	},
	{
		name: "Check worker is OK (processing task error)",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 1, 1)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{0})
			if err != nil {
				return false
			}

			// чтобы успел выполниться worker в горутине
			time.Sleep(time.Second)

			task := scheduler.GetTask(uuid)
			if task.uuid == "" {
				return false
			}
			if task.status != StatusError {
				return false
			}
			if task.response != nil {
				return false
			}

			return true
		},
	},
}

// -----------------------
// For testing only
func (s *Scheduler) isClosed() bool {
	select {
	case _, opened := <-s.taskQueue:
		if !opened {
			return true
		}
		return false
	default:
		return false
	}
}

func (s *Scheduler) isEmptyTaskQueue() bool {
	select {
	case <-s.taskQueue:
		return false
	default:
		return true
	}
}

// mockers

// Processor
type MockProcessor interface {
	Process([]byte) ([]byte, error)
}

type mockprocessor struct{}

func (m *mockprocessor) Process(in []byte) ([]byte, error) {
	if slices.Equal(in, []byte{100}) {
		return []byte{150}, nil
	}

	if slices.Equal(in, []byte{0}) {
		return nil, fmt.Errorf("error processing")
	}

	return in, nil
}

func NewMockProcessor() MockProcessor {
	return &mockprocessor{}
}

func makeProcessor() processor {
	return NewMockProcessor()
}

// Mock Storage

type MockStorage[T any] interface {
	Store(Task) UUID
	Get(UUID) Task
	Find([]FindOperator) []Task
}

type mockstorage[T any] struct {
	tasks map[UUID]Task
}

func (m *mockstorage[any]) Store(t Task) UUID {
	m.tasks[t.uuid] = t
	return t.uuid
}

func (m *mockstorage[any]) Get(uuid UUID) Task {
	if uuid == "" {
		return Task{}
	}
	val, ok := m.tasks[uuid]
	if !ok {
		return Task{}
	}
	return val
}

func (m *mockstorage[any]) Find(operators []FindOperator) []Task {
	foundedTasks := []Task{}

	for _, oper := range operators {
		for _, v := range m.tasks {
			if v.hash == oper.value {
				foundedTasks = append(foundedTasks, v)
			}
		}
	}

	return foundedTasks
}

func NewMockStorage() MockStorage[any] {
	return &mockstorage[any]{
		tasks: make(map[UUID]Task),
	}
}

func makeStorage() storage[Task] {
	return NewMockStorage()
}
