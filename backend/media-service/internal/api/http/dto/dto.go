package dto

import "mime/multipart"

type PostMediaRequest struct {
	OwnerType  string                `form:"owner_type" validate:"required,oneof=coworking"`
	OwnerID    string                `form:"owner_id" validate:"required,uuid4"`
	Purpose    string                `form:"purpose" validate:"required,oneof=cover gallery"`
	UploadedBy string                `form:"uploaded_by" validate:"required,uuid4"`
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
