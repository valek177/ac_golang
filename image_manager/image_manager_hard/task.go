package main

// Сигнатура нашего сервиса
type ImageManagerServiceHandler interface {
	// TODO: реализовать
	UploadImage(url string) (string, error)
}

// Адаптер для взаимодействия с хранилищем картинок
type ImageStorageAdapter interface {
	UploadImage(id string, data []byte) error
	GetImageByID(id string) ([]byte, error)
}

// Адаптер для взаимодействия с БД картинок
type ImageURLDatabaseAdapter interface {
	// TODO только методы
}
