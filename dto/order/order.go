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
	UpdateStatus(req model.UpdateStatus) error
	UpdateShipping(req model.RequestUpdateShipping) (*string, error)
	CreateShippingLogs(orderID uuid.UUID, refCode, shippingStatusFrom, shippingStatusTo, notes string) error
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
		FromStatus: "",
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
		ToStatus:   model.OrderStatusPaid,
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

func (o *order) UpdateStatus(req model.UpdateStatus) (*string, error) {
	successMessage := "Update status success"
	if err := o.store.UpdateStatus(req); err != nil {
		return nil, err
	}

	return &successMessage, nil
}

func (o *order) UpdateShipping(req model.RequestUpdateShipping) (*string, error) {
	successMessage := "Update shipping status success"
	refCode, err := o.store.UpdateShipping(req)
	if err != nil {
		return nil, err
	}

	if err := o.store.CreateShippingLogs(req.OrderID, *refCode, "", req.ShippingStatusFrom, req.ShippingStatusTo); err != nil {
		return nil, err
	}

	return &successMessage, nil
}
