Имеется некоторый процессор, который реализует долгую ресурсоемкую операцию.
Нам нужно реализовать к нему обертку-шедулер, которая будет обрабатывать входящие запросы, не запуская единовременно обработок выше определенного порога.

Предполагаем, что сверху шедулера будет какой-то веб-сервер (например, gRPC), который будет автоматически сгенерирован по публичным методам шедулера.

Сигнатуры методов processor, storage менять нельзя.

1. Easy version
Нужно:
- спроектировать интерфейс и реализовать методы шедулера, который будет позволять обрабатывать внешние запросы, но не превышая потребление ресурсов (ограничиваем число одновременных обработок)
- работаем в cloud-native среде
- для сохранения данных есть быстрое персистентное хранилище, которое описывается интерфейсом. По сути база данных. Хранилище может делать Create/Update, поиск конкретной записи по primary key, и поиска набора записей по критериям.

2. Hard version
Нужно:
- спроектировать интерфейс и реализовать методы шедулера, который будет позволять обрабатывать внешние запросы, но не превышая потребление ресурсов (ограничиваем число одновременных обработок)
- работаем в cloud-native среде
- для сохранения данных есть быстрое персистентное хранилище, которое описывается интерфейсом. По сути база данных. Хранилище может делать Create/Update, поиск конкретной записи по primary key, и поиска набора записей по критериям.
- в базе данных храним хэши записей. При получении нового запроса сначала ищем существующую запись по хэшу, если она есть, возвращаем ее uuid.

Начальный шаблон

```
package scheduler

type processor interface {
    Process([]byte) ([]byte, error)
}

type FindOperator struct {
    key, operator, value string
}

type UUID string

type storage[T any] interface {
    Store(T) UUID 
    Get(UUID) T 
    Find([]FindOperator) []T // select 
}

type Task struct {
    //TODO
}

const (
	StatusProcessing = "processing"
	StatusQueued     = "queued"
	StatusDone       = "done"
	StatusError      = "error"
)

type Scheduler struct {
    //TODO
}

//TODO scheduler methods
func NewScheduler(st storage[Task], proc processor, numWorkers, queueSize int) (*Scheduler, error) {
    //TODO
}

func (s *Scheduler) worker() {
    //TODO
}

func (s *Scheduler) AddTask(request []byte) (UUID, error) {
    //TODO
}

func (s *Scheduler) GetTask(uuid UUID) Task {
    //TODO
}

func (s *Scheduler) Close() {
    //TODO
}
```
