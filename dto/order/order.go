package order

import (
	"cart-order-service/repository/model"

	"github.com/google/uuid"
)

type orderStore interface {
	CreateOrder(bReq model.Order) (*uuid.UUID, *string, error)
	CreateOrderItemsLogs(bReq model.OrderItemsLogs) (*string, error)
	UpdateOrder(bReq model.RequestCallback) (*string, error)
	GetOrderStatus(userID, orderID uuid.UUID) (*model.Order, error)
}

type order struct {
	store orderStore
}

func NewOrder(store orderStore) *order {
	return &order{store}
}

func (o *order) CreateOrder(bReq model.Order) (*uuid.UUID, error) {
	orderID, refCode, err := o.store.CreateOrder(bReq)
	if err != nil {
		return nil, err
	}

	_, err = o.store.CreateOrderItemsLogs(model.OrderItemsLogs{
		OrderID:    *orderID,
		RefCode:    *refCode,
		FromStatus: model.OrderStatusProcessing,
		ToStatus:   model.OrderStatusPending,
		Notes:      "Order created",
	})
	if err != nil {
		return nil, err
	}

	return orderID, nil
}

func (o *order) CallbackPayment(bReq model.RequestCallback) (*string, error) {
	message := "Payment Success"
	refCode, err := o.store.UpdateOrder(bReq)
	if err != nil {
		return nil, err
	}

	_, err = o.store.CreateOrderItemsLogs(model.OrderItemsLogs{
		OrderID:    bReq.OrderID,
		RefCode:    *refCode,
		FromStatus: model.OrderStatusPending,
		ToStatus:   model.OrderStatusCompleted,
		Notes:      "Payment success",
	})
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (o *order) GetOrderStatus(userID, orderID uuid.UUID) (*model.Order, error) {
	return o.store.GetOrderStatus(userID, orderID)
}
