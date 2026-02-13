package dto

type CreateSeverityLevelRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

type UpdateSeverityLevelRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}
