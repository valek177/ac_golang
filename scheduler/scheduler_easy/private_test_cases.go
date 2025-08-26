package main

var privateTestCases = []TestCase{
	{
		name: "Множество задач исполняются параллельно",
		check: func() bool {
			parallelTasksCount := 10
			processorChannel := make(chan struct{}, parallelTasksCount)
			repoChan := make(chan struct{}, parallelTasksCount)

			scheduler, err := NewScheduler(makeRepositoryWithChannel(repoChan),
				makeLongProcessor(processorChannel), parallelTasksCount, 50, generateUUID)
			if err != nil {
				return false
			}
			defer scheduler.Close()

			defer func() {
				for range parallelTasksCount {
					// отправляем таски дальше на обработку в process
					processorChannel <- struct{}{}
				}
			}()

			checkUUIDs := make([]UUID, parallelTasksCount)

			for i := 1; i <= parallelTasksCount; i++ {
				uuid, err := scheduler.AddTask([]byte{byte(i)})
				if uuid == "" || err != nil {
					return false
				}

				checkUUIDs = append(checkUUIDs, uuid)
			}

			for i := 0; i < parallelTasksCount; i++ {
				// нужно дождаться сохранения состояния таски в статусе Processing
				<-repoChan
			}
			countProcessedTasks := 0
			for _, uuid := range checkUUIDs {
				// смотрим состояние тасок
				// должны выполняться только workers count тасок
				task := scheduler.GetTask(uuid)

				if task.status == StatusProcessing {
					countProcessedTasks++
				}

			}

			if countProcessedTasks > parallelTasksCount {
				return false
			}

			return true
		},
	},
	{
		name: "Превышение размера очереди при добавлении задачи (AddTask)",
		check: func() bool {
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(),
				1, 1, generateUUID)
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
			scheduler, err := NewScheduler(makeRepository(), makeProcessor(), 1, 1, generateUUID)
			if err != nil {
				return false
			}

			scheduler.Close()

			return scheduler.isClosed()
		},
	},
}
