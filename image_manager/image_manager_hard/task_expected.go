package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	urllib "net/url"
)

var (
	ErrInvalidURL          = fmt.Errorf("invalid url")
	ErrInternalServer      = fmt.Errorf("internal server error")
	ErrAlreadyUploadingImg = fmt.Errorf("already uploading")
)

type customError struct {
	Message    string
	StatusCode int
}

func (e *customError) Error() string {
	return e.Message
}

// Image статусы
const (
	StatusUploaded  = "uploaded"
	StatusUploading = "uploading"
	StatusError     = "error"
)

type ImageManagerServiceHandler interface {
	UploadImage(url string) (string, error)
}

type ImageManagerService struct {
	adapterStorage ImageStorageAdapter
	adapterDB      ImageURLDatabaseAdapter
}

// Адаптер для взаимодействия с хранилищем картинок
type ImageStorageAdapter interface {
	UploadImage(id string, data []byte) error
	GetImageByID(id string) ([]byte, error)
}

// Адаптер для взаимодействия с бд картинок
type ImageURLDatabaseAdapter interface {
	// TODO реализовать только методы
	GetImageStatus(id string) (string, error)
	UpdateImage(id string, status string) error
}

func NewImageManagerServiceHandler(imageStorageAdapter ImageStorageAdapter,
	adapterDB ImageURLDatabaseAdapter,
) (*ImageManagerService, error) {
	return &ImageManagerService{
		adapterStorage: imageStorageAdapter,
		adapterDB:      adapterDB,
	}, nil
}

func (s *ImageManagerService) UploadImage(url string) (string, error) {
	if !isUrlValid(url) {
		return "", &customError{
			Message:    ErrInvalidURL.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}
	id := generateIdFromUrl(url)

	status, err := s.adapterDB.GetImageStatus(id)
	if err != nil {
		return "", &customError{
			Message:    ErrInternalServer.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	if status == StatusUploaded {
		return id, nil
	}

	if status == StatusUploading {
		return id, &customError{
			Message:    ErrAlreadyUploadingImg.Error(),
			StatusCode: http.StatusConflict,
		}
	}

	err = s.adapterDB.UpdateImage(id, StatusUploading)
	if err != nil {
		return id, &customError{
			Message:    ErrInternalServer.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	data, err := getDataFromURL(url)
	if err != nil {
		s.adapterDB.UpdateImage(id, StatusError)
		return "", &customError{
			Message:    ErrInternalServer.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = s.adapterStorage.UploadImage(id, data)
	if err != nil {
		return "", &customError{
			Message:    ErrInternalServer.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	err = s.adapterDB.UpdateImage(id, StatusUploaded)
	if err != nil {
		return "", &customError{
			Message:    ErrInternalServer.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return id, nil
}

// можем предложить как готовую функцию
func generateIdFromUrl(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))
}

// можем предложить как готовую функцию
func isUrlValid(url string) bool {
	_, err := urllib.ParseRequestURI(url)
	if err != nil {
		return false
	}

	return true
}

// можем предложить как готовую функцию
func getDataFromURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
