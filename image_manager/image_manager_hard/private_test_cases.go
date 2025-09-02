package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	urllib "net/url"
)

var privateTestCases = []TestCase{
	{
		name: "Загрузка картинки с невалидным URL",
		check: func() bool {
			ctx := context.TODO()
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, makeMockURLData())

			_, err = imgManager.UploadImage(ctx, errorURL)
			if err == nil {
				return false
			}

			if err != ErrInvalidURL {
				return false
			}

			return true
		},
	},
	{
		name: "Получение id картинки при попытке загрузить, т.к. она уже загружена",
		check: func() bool {
			ctx := context.TODO()
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, makeMockURLData())

			id, err := imgManager.UploadImage(ctx, uploadedImgURL)
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
			ctx := context.TODO()
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, makeMockURLData())

			id, err := imgManager.UploadImage(ctx, uploadingImgURL)
			if err == nil {
				return false
			}

			if id != generateIdFromUrl(uploadingImgURL) {
				return false
			}

			if err != ErrAlreadyUploadingImg {
				return false
			}

			return true
		},
	},
	{
		name: "Ошибка при попытке получения картинки с URL",
		check: func() bool {
			ctx := context.TODO()
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, makeMockURLData())

			id, err := imgManager.UploadImage(ctx, uploadingImgErrorURL)
			if err == nil {
				return false
			}

			if id != generateIdFromUrl(uploadingImgErrorURL) {
				return false
			}

			if err != ErrInternalServer {
				return false
			}

			return true
		},
	},
	{
		name: "Ошибка при загрузке картинки в хранилище",
		check: func() bool {
			ctx := context.TODO()
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, makeMockURLData())

			id, err := imgManager.UploadImage(ctx, uploadingImgToStorageErrURL)
			if err == nil {
				return false
			}

			if id != generateIdFromUrl(uploadingImgToStorageErrURL) {
				return false
			}

			if err != ErrInternalServer {
				return false
			}

			return true
		},
	},
	{
		name: "Ошибка при обновлении статуса картинки после ее загрузки",
		check: func() bool {
			ctx := context.TODO()
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, makeMockURLData())

			id, err := imgManager.UploadImage(ctx, uploadedImgUpdStatusErrURL)
			if err == nil {
				return false
			}

			if id != generateIdFromUrl(uploadedImgUpdStatusErrURL) {
				return false
			}

			if err != ErrInternalServer {
				return false
			}

			return true
		},
	},
}

func generateIdFromUrl(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))
}

func isUrlValid(url string) bool {
	u, err := urllib.Parse(url)
	return err == nil && u.Scheme != "" && u.Host != ""
}
