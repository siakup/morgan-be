package domain

import "time"

type Division struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Status    bool       `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *string    `json:"created_by"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy *string    `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy *string    `json:"deleted_by"`
}

type DivisionFilter struct {
	Limit  int
	Offset int
	Search string
}
