package main

import (
	"context"
	"sync"
)

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
	{
		name: "Возврат ошибки при загрузке картинки из нескольких потоков",
		check: func() bool {
			adapterStorageChan := make(chan struct{})
			defer close(adapterStorageChan)
			dbChannel := make(chan struct{})
			defer close(dbChannel)

			ctx := context.TODO()
			imgManager, err := NewImageManagerServiceHandler(
				makeImageStorageAdapterWithChannel(adapterStorageChan),
				makeImageURLDatabaseAdapterWithChannel(dbChannel), generateIdFromUrl, makeMockHTTPClient())

			cnt := 10
			wg := sync.WaitGroup{}
			wg.Add(cnt)

			// 1 горутина на основную загрузку картинки
			// отправляем картинку на загрузку, но не загружаем (uploading)
			go func() {
				imgManager.UploadImage(ctx, uploadImgOk)
			}()

			dbChannel <- struct{}{}

			// запускаем загрузку той же картинки в других потоках
			// ожидаем возврат ошибки
			errChan := make(chan error, cnt)
			for i := 0; i < cnt; i++ {
				go func(i int) {
					defer wg.Done()
					_, err = imgManager.UploadImage(ctx, uploadImgOk)
					if err != nil {
						errChan <- err
					}
				}(i)
				dbChannel <- struct{}{}
			}

			go func() {
				wg.Wait()
				close(errChan)
				// завершаем 1 горутину (загрузку картинки)
				adapterStorageChan <- struct{}{}
			}()

			errorsCnt := 0
			for err := range errChan {
				errorsCnt++
				if err == nil || err != ErrAlreadyUploadingImg {
					return false
				}
			}

			if errorsCnt != cnt {
				return false
			}

			return true
		},
	},
}
