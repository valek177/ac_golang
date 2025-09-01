package main

import (
	"context"
	"fmt"
)

type MockImageURLDatabaseAdapter interface {
	UpdateImage(ctx context.Context, id string, status string) error
	PutImage(ctx context.Context, id string, url string) error
	GetImage(ctx context.Context, id string) (url, status string, err error)
}

type mockimageurldatabaseadapter struct{}

func (db *mockimageurldatabaseadapter) GetImage(ctx context.Context, id string) (string, string, error) {
	if id == generateIdFromUrl(errorStatusImg) {
		return "", "", fmt.Errorf("unable to get image status")
	}

	if id == generateIdFromUrl(uploadedImgURL) {
		return "", StatusUploaded, nil
	}

	if id == generateIdFromUrl(uploadingImgURL) {
		return "", StatusUploading, nil
	}

	return "", "", nil
}

func (db *mockimageurldatabaseadapter) UpdateImage(ctx context.Context, id, status string) error {
	if id == generateIdFromUrl(uploadingImgErrorURL) {
		return fmt.Errorf("unable to get image status")
	}
	return nil
}

func (db *mockimageurldatabaseadapter) PutImage(ctx context.Context, id, url string) error {
	return nil
}

func NewMockImageURLDatabaseAdapter() MockImageURLDatabaseAdapter {
	return &mockimageurldatabaseadapter{}
}

func makeImageURLDatabaseAdapter() ImageURLDatabaseAdapter {
	return NewMockImageURLDatabaseAdapter()
}

// MockImageStorageAdapter
type MockImageStorageAdapter interface {
	UploadImage(ctx context.Context, id string, data []byte) error
	GetImageByID(ctx context.Context, id string) ([]byte, error)
}

type mockimagestorageadapter struct{}

func (st *mockimagestorageadapter) UploadImage(ctx context.Context, id string, data []byte) error {
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
