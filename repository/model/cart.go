package model

import (
	"github.com/google/uuid"
	"time"
)

type Cart struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	ProductID uuid.UUID  `json:"product_id"`
	Qty       int        `json:"qty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type GetCartRequest struct {
	UserID    uuid.UUID   `json:"user_id"`
	ProductID []uuid.UUID `json:"product_id"`
}

type DeleteCartRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
}
