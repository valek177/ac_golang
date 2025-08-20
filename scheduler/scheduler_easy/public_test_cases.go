package main

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Check NewScheduler is OK with valid parameters",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 5, 1)
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
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 5, 1)
			if err != nil {
				return false
			}

			task := scheduler.GetTask("123")

			if task.uuid == "" {
				return false
			}

			return true
		},
	},
	{
		name: "Check AddTask is OK",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 5, 1)
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
		name: "Check AddTask is not OK (pool is full)",
		check: func() bool {
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 5, 1)
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
			scheduler, err := NewScheduler(makeStorage(), makeProcessor(), 5, 1)
			if err != nil {
				return false
			}

			scheduler.Close()

			// will occure panic: send on closed channel
			// uuid, err := scheduler.AddTask([]byte{1})
			// if err == nil || uuid != "" {
			// 	return false
			// }

			return scheduler.isClosed()
		},
	},
}

// -----------------------
// For testing only
func (s *Scheduler) isClosed() bool {
	select {
	case _, opened := <-s.closeDoneCh:
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

func (m *mockprocessor) Process([]byte) ([]byte, error) {
	return nil, nil
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

type mockstorage[T any] struct{}

func (m *mockstorage[any]) Store(t Task) UUID {
	return ""
}

func (m *mockstorage[any]) Get(uuid UUID) Task {
	if uuid == "123" {
		return Task{uuid: "123"}
	}
	return Task{}
}

func (m *mockstorage[any]) Find([]FindOperator) []Task {
	return nil
}

func NewMockStorage() MockStorage[any] {
	return &mockstorage[any]{}
}

func makeStorage() storage[Task] {
	return NewMockStorage()
}
