package main

var privateTestCases = []TestCase{
	{
		name: "GetTask возвращает не пустую задачу",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 1, generateUUID)
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
		name: "Добавляем 1 задачу на выполнение AddTask",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 1, generateUUID)
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
		name: "Вызов Close",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 1, generateUUID)
			if err != nil {
				return false
			}

			scheduler.Close()

			return scheduler.isClosed()
		},
	},
}
