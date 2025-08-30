package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	adapterStorage *ImageStorage
	adapterDB      *ImageURLDatabase
}

// Адаптер для взаимодействия с хранилищем картинок
type ImageStorageAdapter interface {
	UploadImage(id string, data []byte) error
	GetImageByID(id string) ([]byte, error)
}

type ImageStorage struct{}

// Адаптер для взаимодействия с бд картинок
type ImageURLDatabaseAdapter interface {
	// TODO реализовать только методы
	GetImageStatus(id string) (string, error)
	UpdateImage(id string, status string)
}

type ImageURLDatabase struct{}

func NewImageManagerServiceHandler(imageStorageAdapter *ImageStorage,
	adapterDB *ImageURLDatabase,
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

	// начинаем транзакцию--->
	err = s.adapterDB.UpdateImage(id, StatusUploading)
	if err != nil {
		return id, &customError{
			Message:    ErrInternalServer.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	data, err := GetDataFromURL(url)
	if err != nil {
		s.adapterDB.UpdateImage(id, StatusError)
		return "", &customError{
			Message:    ErrInternalServer.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	// заканчиваем транзакцию

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

// TODO mockers
func (db *ImageURLDatabase) GetImageStatus(id string) (string, error) {
	return "", nil
}

func (db *ImageURLDatabase) UpdateImage(id string, status string) error {
	return nil
}

func (st *ImageStorage) UploadImage(id string, data []byte) error {
	return nil
}

func (st *ImageStorage) GetImageByID(id string) ([]byte, error) {
	return nil, nil
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

func GetDataFromURL(url string) ([]byte, error) {
	return nil, nil
}

// func main() {
// 	imgStorage := NewSt
// 	service := NewImageManagerServiceHandler()
// 	// Register a handler for the root path
// 	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello from Go's HTTP server!")
// 	})

// 	// Start the server on port 8080
// 	fmt.Println("Server listening on :8085")
// 	log.Fatal(http.ListenAndServe(":8085", nil))
// }
