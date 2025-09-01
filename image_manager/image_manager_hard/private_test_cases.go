package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	urllib "net/url"
)

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
			ctx := context.TODO()
			var urlData URLData // need mocking
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, urlData)

			_, err = imgManager.UploadImage(ctx, errorURL)
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
			ctx := context.TODO()
			var urlData URLData
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, urlData)

			_, err = imgManager.UploadImage(ctx, errorStatusImg)
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
			ctx := context.TODO()
			var urlData URLData
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, urlData)

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
			var urlData URLData
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, urlData)

			id, err := imgManager.UploadImage(ctx, uploadingImgURL)
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
			ctx := context.TODO()
			var urlData URLData
			imgManager, err := NewImageManagerServiceHandler(makeImageStorageAdapter(),
				makeImageURLDatabaseAdapter(), generateIdFromUrl, urlData)

			id, err := imgManager.UploadImage(ctx, uploadingImgErrorURL)
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

func generateIdFromUrl(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))
}

func isUrlValid(url string) bool {
	_, err := urllib.ParseRequestURI(url)
	if err != nil {
		return false
	}

	return true
}

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
