package model

import (
	"time"

	"github.com/google/uuid"
)

var (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"
)

type Order struct {
	UserID        uuid.UUID   `json:"user_id"`
	PaymentTypeID uuid.UUID   `json:"payment_type_id"`
	OrderNumber   string      `json:"order_number"`
	SubtotalPrice float64     `json:"subtotal_price"`
	TotalPrice    float64     `json:"total_price"`
	ProductOrder  interface{} `json:"product_order"`
	Status        string      `json:"status"`
	IsPaid        bool        `json:"is_paid"`
	RefCode       string      `json:"ref_code"`
	CreatedAt     *time.Time  `json:"created_at"`
	UpdatedAt     *time.Time  `json:"updated_at"`
	DeleredAt     *time.Time  `json:"deleted_at"`
}

type OrderItemsLogs struct {
	OrderID    uuid.UUID  `json:"order_id"`
	RefCode    string     `json:"ref_code"`
	FromStatus string     `json:"from_status"`
	ToStatus   string     `json:"to_status"`
	Notes      string     `json:"notes"`
	CreatedAt  *time.Time `json:"created_at"`
}
