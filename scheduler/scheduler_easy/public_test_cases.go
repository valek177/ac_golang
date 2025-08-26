package main

import (
	"slices"
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
			waitingChannel := make(chan struct{})
			scheduler, err := NewScheduler(makeRepository(),
				makeProcessorWithChannel(waitingChannel), 1, 1, generateUUID)
			if err != nil {
				return false
			}

			defer scheduler.Close()

			uuid, err := scheduler.AddTask([]byte{1})
			if err != nil {
				return false
			}

			// ждем, когда таска обработается и попадет в репозиторий
			<-waitingChannel

			task := scheduler.GetTask(uuid)
			if task.uuid == "" {
				return false
			}

			if task.status == StatusDone &&
				!slices.Equal([]byte{1}, task.response) {
				return false
			}

			// таска может не успеть выполниться
			if !slices.Contains([]string{StatusProcessing, StatusDone},
				task.status) {
				return false
			}

			return true
		},
	},
	{
		name: "Проверка работы worker (возврат ошибки при выполнении задачи)",
		check: func() bool {
			waitingChannel := make(chan struct{})
			scheduler, err := NewScheduler(makeRepository(),
				makeProcessorWithChannel(waitingChannel), 1, 1, generateUUID)
			if err != nil {
				return false
			}

			defer scheduler.Close()

			uuid, err := scheduler.AddTask([]byte{0})
			if err != nil {
				return false
			}

			// ждем, когда таска обработается и попадет в репозиторий
			<-waitingChannel

			task := scheduler.GetTask(uuid)
			if task.uuid == "" {
				return false
			}
			if !slices.Contains([]string{StatusProcessing, StatusError},
				task.status) {
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
