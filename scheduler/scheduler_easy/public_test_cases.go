package main

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Check NewScheduler",
		check: func() bool {
			scheduler := NewScheduler(makeStorage(), makeProcessor(), 5, 1)

			if scheduler.st == nil {
				return false
			}

			if scheduler.proc == nil {
				return false
			}

			if scheduler.numWorkers == 0 {
				return false
			}

			if scheduler.taskQueue == nil {
				return false
			}

			return true
		},
	},
	{
		name: "Check GetTask",
		check: func() bool {
			scheduler := NewScheduler(makeStorage(), makeProcessor(), 5, 1)

			task := scheduler.GetTask("123")

			if task.uuid == "" {
				return false
			}

			return true
		},
	},
	{
		name: "Check AddTask",
		check: func() bool {
			scheduler := NewScheduler(makeStorage(), makeProcessor(), 5, 1)

			scheduler.AddTask([]byte{1})
			// todo

			return true
		},
	},
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
