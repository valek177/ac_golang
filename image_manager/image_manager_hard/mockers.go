package main

import "fmt"

type MockImageURLDatabaseAdapter interface {
	GetImageStatus(id string) (string, error)
	UpdateImage(id string, status string) error
}

type mockimageurldatabaseadapter struct{}

func (db *mockimageurldatabaseadapter) GetImageStatus(id string) (string, error) {
	if id == generateIdFromUrl(errorStatusImg) {
		return "", fmt.Errorf("unable to get image status")
	}

	if id == generateIdFromUrl(uploadedImgURL) {
		return StatusUploaded, nil
	}

	if id == generateIdFromUrl(uploadingImgURL) {
		return StatusUploading, nil
	}

	return "", nil
}

func (db *mockimageurldatabaseadapter) UpdateImage(id string, status string) error {
	if id == generateIdFromUrl(uploadingImgErrorURL) {
		return fmt.Errorf("unable to get image status")
	}
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
	UploadImage(id string, data []byte) error
	GetImageByID(id string) ([]byte, error)
}

type mockimagestorageadapter struct{}

func (st *mockimagestorageadapter) UploadImage(id string, data []byte) error {
	return nil
}

func (st *mockimagestorageadapter) GetImageByID(id string) ([]byte, error) {
	return []byte{1}, nil
}

func NewMockImageStorageAdapter() MockImageStorageAdapter {
	return &mockimagestorageadapter{}
}

func makeImageStorageAdapter() ImageStorageAdapter {
	return NewMockImageStorageAdapter()
}
