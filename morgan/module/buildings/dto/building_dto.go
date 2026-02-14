package dto

type CreateBuildingRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

type UpdateBuildingRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}
