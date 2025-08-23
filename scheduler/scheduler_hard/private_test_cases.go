package main

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
		name: "Проверка AddTask",
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

			return !scheduler.isEmptyTaskQueue()
		},
	},
	{
		name: "Проверка Close",
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
