package dto

import "mime/multipart"

type PostMediaRequest struct {
	File       *multipart.FileHeader `form:"file" validate:"required"`
}

type PostMediaResponse struct {
	ID     string            `json:"id"`
	Status string            `json:"status"`
	URLs   map[string]string `json:"urls"`
}

type DeleteMediaRequest struct {
	ID string `param:"id" validate:"required,uuid4"`
}

type ReorderMediaRequest struct {
	Orders map[string]int `json:"orders" validate:"required"`
}
