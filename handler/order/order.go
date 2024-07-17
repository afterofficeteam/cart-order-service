package order

import (
	"cart-order-service/helper/request"
	"cart-order-service/helper/response"
	"cart-order-service/repository/model"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type orderDto interface {
	CreateOrder(bReq model.Order) (*uuid.UUID, error)
}

type Handler struct {
	order     orderDto
	validator *validator.Validate
}

func NewHandler(order orderDto, validator *validator.Validate) *Handler {
	return &Handler{order, validator}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var bReq model.Order
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	bReq.RefCode = request.GenerateRefCode()

	if err := h.validator.Struct(bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	bResp, err := h.order.CreateOrder(bReq)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusCreated, bResp)
}

func (h *Handler) CallbackPayment(w http.ResponseWriter, r *http.Request) {

}
