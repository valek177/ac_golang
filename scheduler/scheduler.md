# Easy version

Есть процессор, который выполняет долгие и ресурсоёмкие операции. Нужно сделать для него **обёртку-шедулер**, которая:
- принимает внешние запросы,
- запускает их обработку через процессор,
- ограничивает число одновременных обработок (не выше заданного порога).

Поверх шедулера будет веб-сервер (например, gRPC), генерируемый автоматически по его публичным методам.

Ограничения:
1. Сигнатуры методов Processor, Repository менять нельзя.
2. Для хранения данных используется быстрое персистентное хранилище (база данных - Repository), которое поддерживает сохранение/обновление, поиск по UUID.

Нужно:
- спроектировать интерфейс шедулера
- реализовать его методы так, чтобы обрабатывать внешние запросы без превышения лимита одновременных операций
- использовать хранилище для работы с данными
- для шедулера нужно иметь возможность прекратить добавлять задачи на исполнение и возвращать ошибку, а оставшиеся задачи в очереди должны быть исполнены

# Hard version

Есть процессор, который выполняет долгие и ресурсоёмкие операции. Нужно сделать для него **обёртку-шедулер**, которая:
- принимает внешние запросы,
- запускает их обработку через процессор,
- ограничивает число одновременных обработок (не выше заданного порога).

Поверх шедулера будет веб-сервер (например, gRPC), генерируемый автоматически по его публичным методам.

Ограничения:
1. Сигнатуры методов Processor, Repository менять нельзя.
2. Для хранения данных используется быстрое персистентное хранилище (база данных - Repository), которое поддерживает сохранение/обновление, поиск по UUID и хэшу запроса.
3. Запросы часто повторяются, поэтому нужно избегать повторной обработки. Результат процессинга всегда актуален и не устаревает, поэтому задачи следует искать в базе данных по хэшу запроса и переиспользовать уже готовый результат.

Нужно:
- спроектировать интерфейс шедулера
- реализовать его методы так, чтобы обрабатывать внешние запросы без превышения лимита одновременных операций
- использовать хранилище для работы с данными
- для шедулера нужно иметь возможность прекратить добавлять задачи на исполнение и возвращать ошибку, а оставшиеся задачи в очереди должны быть исполнены

# Начальный шаблон для easy

```
package scheduler

import "crypto/rand"

type Processor interface {
    Process([]byte) ([]byte, error)
}

type UUID string

type Repository interface {
	Store(Task) UUID
	GetByUUID(UUID) Task
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
func NewScheduler(repository Repository, processor Processor, numWorkers, queueSize int, generateUUID func() UUID) (*Scheduler, error) {
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

// Для генерации UUID можно использовать готовую функцию
func generateUUID() UUID {
	// pseudo uuid
    b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return UUID(uuid)
}
```

# Начальный шаблон для hard

```
package scheduler

import (
	"crypto/sha256"
	"encoding/base64"
)

type Processor interface {
    Process([]byte) ([]byte, error)
}

type (
	UUID string
	Hash string
)

type Repository interface {
	Store(Task) UUID
	GetByUUID(UUID) Task
	GetByHash(hash Hash) []Task
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
func NewScheduler(repository Repository, processor Processor, numWorkers, queueSize int, generateUUID func() UUID, generateHash func(request []byte) Hash) (*Scheduler, error) {
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

// Для генерации UUID можно использовать готовую функцию
func generateUUID() UUID {
	// pseudo uuid
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return UUID(uuid)
}

// Для генерации хэша можно использовать готовую функцию
func generateHash(request []byte) Hash {
	h := sha256.New()

	h.Write(request)

	return Hash(base64.URLEncoding.EncodeToString(h.Sum(nil)))
}
```
