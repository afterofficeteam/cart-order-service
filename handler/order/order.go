package order

import (
	"cart-order-service/helper/response"
	"cart-order-service/repository/model"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type orderDto interface {
	CreateOrder(bReq model.Order) (*uuid.UUID, error)
}

type Handler struct {
	order orderDto
}

func NewHandler(order orderDto) *Handler {
	return &Handler{order}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var bReq model.Order
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if bReq.UserID == uuid.Nil || bReq.PaymentTypeID == uuid.Nil {
		response.HandleResponse(w, http.StatusBadRequest, "user id or payment type id is required")
		return
	}

	if bReq.PaymentTypeID == uuid.Nil {
		response.HandleResponse(w, http.StatusBadRequest, "payment type id is required")
		return
	}

	if bReq.OrderNumber == "" || bReq.Status == "" {
		response.HandleResponse(w, http.StatusBadRequest, "order number or status is required")
		return
	}

	if bReq.SubtotalPrice == 0 {
		response.HandleResponse(w, http.StatusBadRequest, "subtotal price is required")
		return
	}

	if bReq.TotalPrice == 0 {
		response.HandleResponse(w, http.StatusBadRequest, "total price is required")
		return
	}

	if bReq.Status == "" {
		response.HandleResponse(w, http.StatusBadRequest, "status is required")
		return
	}

	orderID, err := h.order.CreateOrder(bReq)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusCreated, orderID)
}
