package main

// Обработчик задач
type processor interface {
	Process([]byte) ([]byte, error)
}

// Для хранилища
type FindOperator struct {
	key, operator, value string
}

type UUID string

// Интерфейс хранилища (БД)
type storage[T any] interface {
	Store(T) UUID
	Get(UUID) T
	Find([]FindOperator) []T
}

// TODO для К.
type Task struct {
	uuid     UUID
	status   string
	request  []byte
	response []byte
}

const (
	StatusProcessing = "processing"
	StatusQueued     = "queued"
	StatusDone       = "done"
	StatusError      = "error"
)

// TODO для К.
type Scheduler struct {
	st   storage[Task]
	proc processor

	taskQueue  chan Task
	numWorkers int
}

func NewScheduler(st storage[Task], proc processor, numWorkers, queueSize int) *Scheduler {
	scheduler := &Scheduler{
		st:         st,
		proc:       proc,
		numWorkers: numWorkers,
		taskQueue:  make(chan Task, queueSize),
	}

	for range numWorkers {
		go scheduler.worker()
	}

	return scheduler
}

func (s *Scheduler) worker() {
	var err error

	for t := range s.taskQueue {
		t.status = StatusProcessing

		s.st.Store(t)

		t.response, err = s.proc.Process(t.request)
		if err != nil {
			t.status = StatusError
		} else {
			t.status = StatusDone
		}

		s.st.Store(t)
	}
}

func (s *Scheduler) AddTask(request []byte) UUID {
	t := Task{
		uuid:    newUUID(),
		status:  StatusQueued,
		request: request,
	}

	go func() {
		s.taskQueue <- t
	}()

	return t.uuid
}

func (s *Scheduler) GetTask(uuid UUID) Task {
	task := s.st.Get(uuid)

	if task.uuid == "" {
		return Task{}
	}

	return task
}

// Генератор UUID
func newUUID() UUID {
	return "d97976cc-35f8-44cb-91f9-fa47a85db34b"
}
