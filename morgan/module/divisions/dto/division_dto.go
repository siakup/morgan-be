package dto

type CreateDivisionRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

type UpdateDivisionRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}
