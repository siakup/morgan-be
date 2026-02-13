package dto

type CreateShiftGroupRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

type UpdateShiftGroupRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}
