На вход веб-сервису приходит HTTP-запрос. Необходимо реализовать production-ready менеджер для заливки картинок по URL в уже готовый сервис для хранения картинок (ImageManagerServiceHandler).

Требования:
1. По данному пользователем URL скопировать изображение во внутренний сервис для хранения картинок и получить внутренний id в сервисе хранения, либо вернуть уже существующий id, если этот URL был загружен ранее.
2. Используем готовый интерфейс клиента для заливки изображения в сервис хранения по URL (ImageStorageAdapter)
3. Для исключения повторной заливки картинки использовать вспомогательную базу данных. Нужно предложить интерфейс для работы с этой базой данных без реализации (ImageURLDatabaseAdapter) и, опираясь на него, реализовать логику в менеджере для обработки запросов.

Требования по обработке ошибок:
- если загрузка не удалась из-за некорректного URL, хотим возвращать 400 (bad request) - ErrInvalidURL
- в случае ошибок, связанных с конфликтующими параллельными запросами, хотим возвращать 409 (conflict)
- в случае ошибок обработки на стороне сервера, хотим возвращать 500 (internal server error)

Допустимо использовать библиотеки:
- "net/url" - для валидации URL


# Начальный шаблон

```
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

type URLData interface {
	Get(url string) ([]byte, error)
}

func NewImageManagerServiceHandler(imageStorageAdapter ImageStorageAdapter,
	adapterDB ImageURLDatabaseAdapter, generateIdFromURL func(url string) string,
	urlData URLData,
) (ImageManagerServiceHandler, error) {
	// TODO
}

```
