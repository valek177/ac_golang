package main

import (
	"context"
	"fmt"
	"net/http"
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

type URLData interface {
	Get(url string) http.Response
	GetBody(response http.Response) []byte
}

type ImageManagerServiceHandler interface {
	UploadImage(url string) (string, error)
}

type ImageManagerService struct {
	adapterStorage    ImageStorageAdapter
	adapterDB         ImageURLDatabaseAdapter
	generateIdFromURL func(url string) string
	urlData           URLData
}

// Адаптер для взаимодействия с хранилищем картинок
type ImageStorageAdapter interface {
	UploadImage(ctx context.Context, id string, data []byte) error
	GetImageByID(ctx context.Context, id string) ([]byte, error)
}

// Адаптер для взаимодействия с бд картинок
type ImageURLDatabaseAdapter interface {
	// TODO реализовать только методы
	UpdateImage(ctx context.Context, id string, status string) error
	PutImage(ctx context.Context, id string, url string) error
	GetImage(ctx context.Context, id string) (url, status string, err error)
}

func NewImageManagerServiceHandler(imageStorageAdapter ImageStorageAdapter,
	adapterDB ImageURLDatabaseAdapter, generateIdFromURL func(url string) string,
	urlData URLData,
) (*ImageManagerService, error) {
	return &ImageManagerService{
		adapterStorage:    imageStorageAdapter,
		adapterDB:         adapterDB,
		generateIdFromURL: generateIdFromURL,
		urlData:           urlData,
	}, nil
}

func (s *ImageManagerService) UploadImage(ctx context.Context, url string) (string, error) {
	if !isUrlValid(url) {
		return "", ErrInvalidURL
	}
	id := generateIdFromUrl(url)

	err := s.adapterDB.PutImage(ctx, id, url)
	if err != nil {
		_, status, err := s.adapterDB.GetImage(ctx, id)
		switch status {
		case StatusUploaded:
			return id, nil
		case StatusUploading:
			return id, ErrAlreadyUploadingImg
		}

		if err != nil {
			return id, ErrInternalServer
		}
	}

	data, err := getDataFromURL(url)
	if err != nil {
		return "", ErrInternalServer
	}

	err = s.adapterStorage.UploadImage(ctx, id, data)
	if err != nil {
		return "", ErrInternalServer
	}

	err = s.adapterDB.UpdateImage(ctx, id, StatusUploaded)
	if err != nil {
		return "", ErrInternalServer
	}

	return id, nil
}
