package main

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Вызов NewScheduler, добавление задачи, получение этой же задачи, проверка статуса задачи",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 1,
				generateUUID, generateHash)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{1})
			if err != nil {
				return false
			}

			scheduler.Close()

			// ждем когда таска обработается, иначе не успеет попасть в repository
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
		name: "Проверка работы worker (возврат ошибки при выполнении задачи)",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 1,
				generateUUID, generateHash)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{0})
			if err != nil {
				return false
			}

			scheduler.Close()
			// ждем когда таска обработается, иначе не успеет попасть в repository
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
	{
		name: "Вызов NewScheduler, проверка коллизии задач (одинаковый хэш)",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 2,
				generateUUID, generateHash)
			if err != nil {
				return false
			}

			checkUuid, err := scheduler.AddTask([]byte{1})
			if err != nil {
				return false
			}

			// ждем когда таска обработается, иначе не успеет попасть в repository
			time.Sleep(time.Second)

			// Добавляем ту же таску, должны вернуть тот же UUID, что у таски в repository
			uuid, err := scheduler.AddTask([]byte{1})
			if uuid != checkUuid || err != nil {
				return false
			}
			scheduler.Close()

			return true
		},
	},
	{
		name: "Вызов NewScheduler, без коллизий задач",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 2,
				generateUUID, generateHash)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{1})
			if err != nil {
				return false
			}

			uuid2, err := scheduler.AddTask([]byte{2})
			if err != nil {
				return false
			}

			scheduler.Close()

			// ждем когда таска обработается, иначе не успеет попасть в repository
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

			task2 := scheduler.GetTask(uuid2)
			if task2.uuid == "" {
				return false
			}

			if !slices.Equal([]byte{2}, task2.response) {
				return false
			}

			if task2.status != StatusDone {
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

func makeProcessor() Processor {
	return NewMockProcessor()
}

// Mock Repository

type MockRepository interface {
	Store(Task) UUID
	GetByUUID(UUID) Task
	GetByHash(hash Hash) []Task
}

type mockrepository struct {
	mutexTasks sync.RWMutex
	tasks      map[UUID]Task
}

func (m *mockrepository) Store(t Task) UUID {
	m.mutexTasks.Lock()
	defer m.mutexTasks.Unlock()

	m.tasks[t.uuid] = t

	return t.uuid
}

func (m *mockrepository) GetByHash(hash Hash) []Task {
	foundedTasks := []Task{}

	for _, v := range m.tasks {
		if v.hash == hash {
			foundedTasks = append(foundedTasks, v)
		}
	}

	return foundedTasks
}

func (m *mockrepository) GetByUUID(uuid UUID) Task {
	if uuid == "" {
		return Task{}
	}
	m.mutexTasks.RLock()
	val, ok := m.tasks[uuid]
	m.mutexTasks.RUnlock()
	if !ok {
		return Task{}
	}
	return val
}

func NewMockRepository() MockRepository {
	return &mockrepository{
		tasks: make(map[UUID]Task),
	}
}

func makeRepository() Repository {
	return NewMockRepository()
}
