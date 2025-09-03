package main

import "context"

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Успешная загрузка картинки",
		check: func() bool {
			ctx := context.TODO()
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, makeMockHTTPClient())

			id, err := imgManager.UploadImage(ctx, uploadImgOk)
			if err != nil {
				return false
			}

			if id != generateIdFromUrl(uploadImgOk) {
				return false
			}

			return true
		},
	},
}
