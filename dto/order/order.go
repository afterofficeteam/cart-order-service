package order

import (
	"cart-order-service/repository/model"

	"github.com/google/uuid"
)

type orderStore interface {
	CreateOrder(bReq model.Order) (*uuid.UUID, *string, error)
	CreateOrderItemsLogs(bReq model.OrderItemsLogs) (*string, error)
}

type order struct {
	store orderStore
}

func NewOrder(store orderStore) *order {
	return &order{store}
}

func (o *order) CreateOrder(bReq model.Order) (*string, error) {
	orderID, refCode, err := o.store.CreateOrder(bReq)
	if err != nil {
		return nil, err
	}

	refCode, err = o.store.CreateOrderItemsLogs(model.OrderItemsLogs{
		OrderID:    *orderID,
		RefCode:    *refCode,
		FromStatus: model.OrderStatusPending,
		ToStatus:   model.OrderStatusProcessing,
		Notes:      "Order created",
	})
	if err != nil {
		return nil, err
	}

	return refCode, nil
}