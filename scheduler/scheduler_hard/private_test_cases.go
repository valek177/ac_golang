package main

import (
	"fmt"
	"slices"
	"time"
)

var privateTestCases = []TestCase{
	{
		name: "GetTask возвращает непустую таску",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 2,
				generateUUID, generateHash)
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
		name: "2 задачи исполняются параллельно",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeLongProcessor(), 2, 2,
				generateUUID, generateHash)
			if err != nil {
				return false
			}
			defer scheduler.Close()

			uuid, err := scheduler.AddTask([]byte{1})
			if uuid == "" || err != nil {
				return false
			}

			uuid2, err := scheduler.AddTask([]byte{2})
			if uuid2 == "" || err != nil {
				return false
			}

			// Чтобы таски успели поместиться в workerы,
			// но не успели выполниться (статус Processing)
			time.Sleep(time.Millisecond * 100)

			task := scheduler.GetTask(uuid)
			if task.status != StatusProcessing {
				return false
			}

			task2 := scheduler.GetTask(uuid2)
			if task2.status != StatusProcessing {
				return false
			}

			return true
		},
	},
	{
		name: "Превышение размера очереди при добавлении задачи (AddTask)",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 1,
				generateUUID, generateHash)
			if err != nil {
				return false
			}

			uuid, err := scheduler.AddTask([]byte{1})
			if uuid == "" || err != nil {
				return false
			}
			// предполагаем, что добавление задачи в очередь происходит практически сразу
			// при этом первая таска не успевает выполниться, и лежит в очереди
			uuid, err = scheduler.AddTask([]byte{2})
			if uuid == "" || err != nil {
				return true
			}

			// если успеет выполниться 1ая таска и не успеет выполнится предыдущая
			uuid, err = scheduler.AddTask([]byte{3})
			if uuid == "" || err != nil {
				return true
			}

			uuid, err = scheduler.AddTask([]byte{4})
			if uuid == "" || err != nil {
				return true
			}

			return false
		},
	},
	{
		name: "Вызов Close",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 2,
				generateUUID, generateHash)
			if err != nil {
				return false
			}

			scheduler.Close()

			return scheduler.isClosed()
		},
	},
}

// Mock Long Processor
type MockLongProcessor interface {
	Process([]byte) ([]byte, error)
}

type mocklongprocessor struct{}

func (m *mocklongprocessor) Process(in []byte) ([]byte, error) {
	// simulate long processing
	time.Sleep(time.Second)

	if slices.Equal(in, []byte{100}) {
		return []byte{150}, nil
	}

	if slices.Equal(in, []byte{0}) {
		return nil, fmt.Errorf("error processing")
	}

	return in, nil
}

func NewMockLongProcessor() MockLongProcessor {
	return &mocklongprocessor{}
}

func makeLongProcessor() Processor {
	return NewMockLongProcessor()
}
