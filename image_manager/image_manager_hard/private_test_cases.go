package main

const (
	errorURL                    = "invalid_url"
	errorStatusImg              = "http://example.com/error_status_img.jpg"
	uploadedImgURL              = "http://example.com/uploaded_img.jpg"
	uploadingImgURL             = "http://example.com/uploading_img.jpg"
	uploadingImgErrorURL        = "http://example.com/error_uploading_img.jpg"
	downloadingImgErrorURL      = "http://example.com/error_downloading_from_url.jpg"
	uploadingImgToStorageErrURL = "http://example.com/error_uploading_to_storage_url.jpg"
)

var privateTestCases = []TestCase{
	{
		name: "Загрузка картинки с невалидным URL",
		check: func() bool {
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter())

			_, err = imgManager.UploadImage(errorURL)
			if err == nil {
				return false
			}

			if err.Error() != ErrInvalidURL.Error() {
				return false
			}

			return true
		},
	},
	{
		name: "Ошибка получения статуса картинки из БД перед ее загрузкой",
		check: func() bool {
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter())

			_, err = imgManager.UploadImage(errorStatusImg)
			if err == nil {
				return false
			}

			if err.Error() != ErrInternalServer.Error() {
				return false
			}

			return true
		},
	},
	{
		name: "Получение id картинки при попытке загрузить, т.к. она уже загружена",
		check: func() bool {
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter())

			id, err := imgManager.UploadImage(uploadedImgURL)
			if err != nil {
				return false
			}

			if id != generateIdFromUrl(uploadedImgURL) {
				return false
			}

			return true
		},
	},
	{
		name: "Получение id картинки и ошибки при попытке загрузить, т.к. она уже загружается",
		check: func() bool {
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter())

			id, err := imgManager.UploadImage(uploadingImgURL)
			if err == nil {
				return false
			}

			if id != generateIdFromUrl(uploadingImgURL) {
				return false
			}

			if err.Error() != ErrAlreadyUploadingImg.Error() {
				return false
			}

			return true
		},
	},
	{
		name: "Ошибка при попытке обновления статуса картинки при ее загрузке",
		check: func() bool {
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter())

			id, err := imgManager.UploadImage(uploadingImgErrorURL)
			if err == nil {
				return false
			}

			if id != generateIdFromUrl(uploadingImgErrorURL) {
				return false
			}

			if err.Error() != ErrInternalServer.Error() {
				return false
			}

			return true
		},
	},
	// {
	// 	name: "Ошибка при загрузке картинки в хранилище",
	// 	check: func() bool {
	// TODO fixme
	// 		imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
	// 			makeImageURLDatabaseAdapter())

	// 		id, err := imgManager.UploadImage(uploadingImgToStorageErrURL)
	// 		if err == nil {
	// 			return false
	// 		}

	// 		if id != generateIdFromUrl(uploadingImgToStorageErrURL) {
	// 			return false
	// 		}

	// 		if err.Error() != ErrInternalServer.Error() {
	// 			return false
	// 		}

	// 		return true
	// 	},
	// },
}
