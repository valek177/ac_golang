package main

import (
	"context"
	"fmt"
	"sync"
)

type MockImageURLDatabaseAdapter interface {
	UpdateImage(ctx context.Context, id string, status string) error
	PutImage(ctx context.Context, id string, url string) error
	GetImage(ctx context.Context, id string) (url, status string, err error)
}

type mockimageurldatabaseadapter struct {
	mutexImages sync.RWMutex
	images      map[string]map[string]string // id: { url: "...", status: "..."}
}

func (db *mockimageurldatabaseadapter) GetImage(ctx context.Context, id string) (url string,
	status string, err error,
) {
	if id == generateIdFromUrl(errorStatusImg) {
		return "", "", fmt.Errorf("unable to get image status")
	}

	if id == generateIdFromUrl(uploadedImgURL) {
		return "", StatusUploaded, nil
	}

	if id == generateIdFromUrl(uploadingImgURL) {
		return "", StatusUploading, nil
	}

	db.mutexImages.RLock()
	defer db.mutexImages.RUnlock()

	val, ok := db.images[id]
	if !ok {
		return "", "", fmt.Errorf("unable to find image with id %s", id)
	}

	return val["url"], val["status"], nil
}

func (db *mockimageurldatabaseadapter) UpdateImage(ctx context.Context, id, status string) error {
	if id == generateIdFromUrl(uploadingImgErrorURL) {
		return fmt.Errorf("unable to get image status")
	} else if id == generateIdFromUrl(uploadedImgUpdStatusErrURL) {
		return fmt.Errorf("unable to update image status")
	}

	db.mutexImages.Lock()
	defer db.mutexImages.Unlock()

	_, ok := db.images[id]
	if !ok {
		return fmt.Errorf("unable to update image with id %s", id)
	}
	db.images[id]["status"] = status

	return nil
}

func (db *mockimageurldatabaseadapter) PutImage(ctx context.Context, id, url string) error {
	if url == uploadingImgURL || url == uploadedImgURL {
		return fmt.Errorf("already exists in database")
	}

	db.mutexImages.Lock()
	defer db.mutexImages.Unlock()
	_, ok := db.images[id]
	if ok {
		return fmt.Errorf("image with id %s already exists", id)
	}

	db.images[id] = map[string]string{
		"url":    url,
		"status": "uploading",
	}

	return nil
}

func NewMockImageURLDatabaseAdapter() MockImageURLDatabaseAdapter {
	return &mockimageurldatabaseadapter{
		images: make(map[string]map[string]string),
	}
}

func makeImageURLDatabaseAdapter() ImageURLDatabaseAdapter {
	return NewMockImageURLDatabaseAdapter()
}

// MockImageURLDatabaseAdapterWithChannel
type MockImageURLDatabaseAdapterWithChannel interface {
	UpdateImage(ctx context.Context, id string, status string) error
	PutImage(ctx context.Context, id string, url string) error
	GetImage(ctx context.Context, id string) (url, status string, err error)
}

type mockimageurldatabaseadapterwithchannel struct {
	mutexImages sync.RWMutex
	images      map[string]map[string]string // id: { url: "...", status: "..."}
	inChannel   chan struct{}
}

func (db *mockimageurldatabaseadapterwithchannel) GetImage(
	ctx context.Context, id string) (url string,
	status string, err error,
) {
	if id == generateIdFromUrl(errorStatusImg) {
		return "", "", fmt.Errorf("unable to get image status")
	}

	if id == generateIdFromUrl(uploadedImgURL) {
		return "", StatusUploaded, nil
	}

	if id == generateIdFromUrl(uploadingImgURL) {
		return "", StatusUploading, nil
	}

	db.mutexImages.RLock()
	defer db.mutexImages.RUnlock()

	val, ok := db.images[id]
	if !ok {
		return "", "", fmt.Errorf("unable to find image with id %s", id)
	}

	return val["url"], val["status"], nil
}

func (db *mockimageurldatabaseadapterwithchannel) UpdateImage(ctx context.Context, id, status string) error {
	if id == generateIdFromUrl(uploadingImgErrorURL) {
		return fmt.Errorf("unable to get image status")
	} else if id == generateIdFromUrl(uploadedImgUpdStatusErrURL) {
		return fmt.Errorf("unable to update image status")
	}

	db.mutexImages.Lock()
	defer db.mutexImages.Unlock()

	_, ok := db.images[id]
	if !ok {
		return fmt.Errorf("unable to update image with id %s", id)
	}
	db.images[id]["status"] = status

	return nil
}

func (db *mockimageurldatabaseadapterwithchannel) PutImage(ctx context.Context, id, url string) error {
	if url == uploadingImgURL || url == uploadedImgURL {
		return fmt.Errorf("already exists in database")
	}

	<-db.inChannel

	db.mutexImages.Lock()
	defer db.mutexImages.Unlock()
	_, ok := db.images[id]
	if ok {
		return fmt.Errorf("image with id %s already exists", id)
	}
	db.images[id] = map[string]string{
		"url":    url,
		"status": "uploading",
	}

	return nil
}

func NewMockImageURLDatabaseAdapterWithChannel(inChannel chan struct{}) MockImageURLDatabaseAdapter {
	return &mockimageurldatabaseadapterwithchannel{
		images:    make(map[string]map[string]string),
		inChannel: inChannel,
	}
}

func makeImageURLDatabaseAdapterWithChannel(inChannel chan struct{},
) ImageURLDatabaseAdapter {
	return NewMockImageURLDatabaseAdapterWithChannel(inChannel)
}

// MockImageStorageAdapter
type MockImageStorageAdapter interface {
	UploadImage(ctx context.Context, id string, data []byte) error
	GetImageByID(ctx context.Context, id string) ([]byte, error)
}

type mockimagestorageadapter struct{}

func (st *mockimagestorageadapter) UploadImage(ctx context.Context, id string, data []byte) error {
	if id == generateIdFromUrl(uploadingImgToStorageErrURL) {
		return fmt.Errorf("unable to upload image")
	}
	return nil
}

func (st *mockimagestorageadapter) GetImageByID(ctx context.Context, id string) ([]byte, error) {
	return []byte{1}, nil
}

func NewMockImageStorageAdapter() MockImageStorageAdapter {
	return &mockimagestorageadapter{}
}

func makeImageStorageAdapter() ImageStorageAdapter {
	return NewMockImageStorageAdapter()
}

// MockImageStorageAdapterWithChannel
type MockImageStorageAdapterWithChannel interface {
	UploadImage(ctx context.Context, id string, data []byte) error
	GetImageByID(ctx context.Context, id string) ([]byte, error)
}

type mockimagestorageadapterwithchannel struct {
	inChannel chan struct{}
}

func (st *mockimagestorageadapterwithchannel) UploadImage(ctx context.Context, id string, data []byte) error {
	if id == generateIdFromUrl(uploadingImgToStorageErrURL) {
		return fmt.Errorf("unable to upload image")
	}

	<-st.inChannel

	return nil
}

func (st *mockimagestorageadapterwithchannel) GetImageByID(ctx context.Context, id string) ([]byte, error) {
	return []byte{1}, nil
}

func NewMockImageStorageAdapterWithChannel(inChannel chan struct{}) MockImageStorageAdapterWithChannel {
	return &mockimagestorageadapterwithchannel{
		inChannel: inChannel,
	}
}

func makeImageStorageAdapterWithChannel(inChannel chan struct{}) ImageStorageAdapter {
	return NewMockImageStorageAdapterWithChannel(inChannel)
}

// HTTPClient

type MockHTTPClient interface {
	Get(url string) ([]byte, error)
}

type mockhttpclient struct{}

func (m *mockhttpclient) Get(url string) ([]byte, error) {
	if url == uploadingImgErrorURL {
		return nil, fmt.Errorf("unable to get image")
	}
	return []byte{1}, nil
}

func NewMockHTTPClient() MockHTTPClient {
	return &mockhttpclient{}
}

func makeMockHTTPClient() HTTPClient {
	return NewMockHTTPClient()
}
