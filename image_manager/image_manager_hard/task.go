package main

import (
	"context"
	"fmt"
)

var (
	ErrInvalidURL          = fmt.Errorf("invalid url")
	ErrInternalServer      = fmt.Errorf("internal server error")
	ErrAlreadyUploadingImg = fmt.Errorf("already uploading")
)

// Image статусы
const (
	StatusUploaded  = "uploaded"
	StatusUploading = "uploading"
	StatusError     = "error"
)

// Сигнатура нашего сервиса
type ImageManagerServiceHandler interface {
	// TODO: реализовать
	UploadImage(ctx context.Context, url string) (string, error)
}

// Адаптер для взаимодействия с хранилищем картинок
type ImageStorageAdapter interface {
	UploadImage(ctx context.Context, id string, data []byte) error
	GetImageByID(ctx context.Context, id string) ([]byte, error)
}

// Адаптер для взаимодействия с БД картинок
type ImageURLDatabaseAdapter interface {
	// TODO только методы
}

type HTTPClient interface {
	Get(url string) ([]byte, error)
}

func NewImageManagerServiceHandler(imageStorageAdapter ImageStorageAdapter,
	adapterDB ImageURLDatabaseAdapter, generateIdFromURL func(url string) string,
	httpClient HTTPClient,
) (ImageManagerServiceHandler, error) {
	// TODO
}
