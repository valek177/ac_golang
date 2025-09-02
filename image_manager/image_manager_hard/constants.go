package main

const (
	uploadImgOk                 = "http://localhost:8080/img.jpg"
	errorURL                    = "http:::/not.valid/a//a??a?b=&&c#hi"
	errorStatusImg              = "http://localhost:8080/error_status_img.jpg"
	uploadedImgURL              = "http://localhost:8080/uploaded_img.jpg"
	uploadingImgURL             = "http://localhost:8080/uploading_img.jpg"
	uploadingImgErrorURL        = "http://localhost:8080/error_uploading_img.jpg"
	downloadingImgErrorURL      = "http://localhost:8080/error_downloading_from_url.jpg"
	uploadingImgToStorageErrURL = "http://localhost:8080/error_uploading_to_storage_url.jpg"
	uploadedImgUpdStatusErrURL  = "http://localhost:8080/error_uploaded_image_upd_status_url.jpg"
)
