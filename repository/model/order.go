package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

var (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"
	OrderStatusPakcing    = "packing"
)

type Order struct {
	UserID        uuid.UUID       `json:"user_id" validate:"required"`
	PaymentTypeID uuid.UUID       `json:"payment_type_id" validate:"required"`
	OrderNumber   string          `json:"order_number" validate:"required"`
	TotalPrice    float64         `json:"total_price" validate:"required"`
	ProductOrder  json.RawMessage `json:"product_order"`
	Status        string          `json:"status" validate:"required"`
	IsPaid        bool            `json:"is_paid"`
	RefCode       string          `json:"ref_code"`
	CreatedAt     *time.Time      `json:"created_at"`
	UpdatedAt     *time.Time      `json:"updated_at"`
	DeletedAt     *time.Time      `json:"deleted_at"`
}

type ProductOrder struct {
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Price         float64   `json:"price"`
	Qty           int       `json:"qty"`
	SubtotalPrice float64   `json:"subtotal_price"`
}

type OrderItemsLogs struct {
	OrderID    uuid.UUID  `json:"order_id"`
	RefCode    string     `json:"ref_code"`
	FromStatus string     `json:"from_status"`
	ToStatus   string     `json:"to_status"`
	Notes      string     `json:"notes"`
	CreatedAt  *time.Time `json:"created_at"`
}

type RequestCallback struct {
	OrderID   uuid.UUID  `json:"order_id" validate:"required"`
	Status    string     `json:"status" validate:"required"`
	IsPaid    bool       `json:"is_paid"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UpdateStatus struct {
	UserID  uuid.UUID `json:"user_id" validate:"required"`
	OrderID uuid.UUID `json:"order_id" validate:"required"`
	Status  string    `json:"status" validate:"required"`
}
