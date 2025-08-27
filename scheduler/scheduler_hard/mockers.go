package main

import (
	"fmt"
	"slices"
	"sync"
)

// mockers

// Processor
type MockProcessor interface {
	Process([]byte) ([]byte, error)
}

type mockprocessor struct{}

func (m *mockprocessor) Process(in []byte) ([]byte, error) {
	if slices.Equal(in, []byte{100}) {
		return []byte{150}, nil
	}

	if slices.Equal(in, []byte{0}) {
		return nil, fmt.Errorf("error processing")
	}

	return in, nil
}

func NewMockProcessor() MockProcessor {
	return &mockprocessor{}
}

func makeProcessor() Processor {
	return NewMockProcessor()
}

// Processor
type MockProcessorWithChannel interface {
	Process([]byte) ([]byte, error)
}

type mockprocessorWithChannel struct {
	waitingChannel chan struct{}
}

func (m *mockprocessorWithChannel) Process(in []byte) ([]byte, error) {
	defer func() {
		m.waitingChannel <- struct{}{}
	}()

	if slices.Equal(in, []byte{100}) {
		return []byte{150}, nil
	}

	if slices.Equal(in, []byte{0}) {
		return nil, fmt.Errorf("error processing")
	}

	return in, nil
}

func NewMockProcessorWithChannel(waitingChannel chan struct{}) MockProcessorWithChannel {
	return &mockprocessorWithChannel{
		waitingChannel: waitingChannel,
	}
}

func makeProcessorWithChannel(waitingChannel chan struct{}) Processor {
	return NewMockProcessorWithChannel(waitingChannel)
}

// Mock Repository

type MockRepository interface {
	Store(Task) UUID
	GetByUUID(UUID) Task
	GetByHash(hash Hash) []Task
}

type mockrepository struct {
	mutexTasks sync.RWMutex
	tasks      map[UUID]Task
}

func (m *mockrepository) Store(t Task) UUID {
	m.mutexTasks.Lock()
	defer m.mutexTasks.Unlock()

	m.tasks[t.uuid] = t

	return t.uuid
}

func (m *mockrepository) GetByHash(hash Hash) []Task {
	foundedTasks := []Task{}

	m.mutexTasks.Lock()
	defer m.mutexTasks.Unlock()
	for _, v := range m.tasks {
		if v.hash == hash {
			foundedTasks = append(foundedTasks, v)
		}
	}

	return foundedTasks
}

func (m *mockrepository) GetByUUID(uuid UUID) Task {
	if uuid == "" {
		return Task{}
	}
	m.mutexTasks.RLock()
	val, ok := m.tasks[uuid]
	m.mutexTasks.RUnlock()

	if !ok {
		return Task{}
	}
	return val
}

func NewMockRepository() MockRepository {
	return &mockrepository{
		tasks: make(map[UUID]Task),
	}
}

func makeRepository() Repository {
	return NewMockRepository()
}

// Mock Long Processor
type MockLongProcessor interface {
	Process([]byte) ([]byte, error)
}

type mocklongprocessor struct {
	startProcess chan struct{}
}

func (m *mocklongprocessor) Process(in []byte) ([]byte, error) {
	<-m.startProcess

	if slices.Equal(in, []byte{100}) {
		return []byte{150}, nil
	}

	if slices.Equal(in, []byte{0}) {
		return nil, fmt.Errorf("error processing")
	}

	return in, nil
}

func NewMockLongProcessor(startProcess chan struct{}) MockLongProcessor {
	return &mocklongprocessor{startProcess: startProcess}
}

func makeLongProcessor(startProcess chan struct{}) Processor {
	return NewMockLongProcessor(startProcess)
}

// Mock Repository

type MockRepositoryWithChannel interface {
	Store(Task) UUID
	GetByUUID(UUID) Task
	GetByHash(hash Hash) []Task
}

type mockrepositoryWithChannel struct {
	mutexTasks sync.RWMutex
	tasks      map[UUID]Task

	in chan struct{}
}

func (m *mockrepositoryWithChannel) Store(t Task) UUID {
	m.mutexTasks.Lock()
	m.tasks[t.uuid] = t
	m.mutexTasks.Unlock()

	if t.status == StatusProcessing {
		m.in <- struct{}{}
	}

	return t.uuid
}

func (m *mockrepositoryWithChannel) GetByHash(hash Hash) []Task {
	foundedTasks := []Task{}

	m.mutexTasks.Lock()
	defer m.mutexTasks.Unlock()
	for _, v := range m.tasks {
		if v.hash == hash {
			foundedTasks = append(foundedTasks, v)
		}
	}

	return foundedTasks
}

func (m *mockrepositoryWithChannel) GetByUUID(uuid UUID) Task {
	if uuid == "" {
		return Task{}
	}
	m.mutexTasks.RLock()
	val, ok := m.tasks[uuid]
	m.mutexTasks.RUnlock()

	if !ok {
		return Task{}
	}

	return val
}

func NewMockRepositoryWithChannel(in chan struct{}) MockRepositoryWithChannel {
	return &mockrepositoryWithChannel{
		tasks: make(map[UUID]Task),
		in:    in,
	}
}

func makeRepositoryWithChannel(in chan struct{}) Repository {
	return NewMockRepositoryWithChannel(in)
}
