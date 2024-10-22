package order

import (
	"cart-order-service/helper/request"
	"cart-order-service/helper/response"
	"cart-order-service/repository/model"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type orderDto interface {
	CreateOrder(bReq model.Order) (*uuid.UUID, error)
	CallbackPayment(bReq model.RequestCallback) (*string, error)
	GetOrderStatus(userID, orderID uuid.UUID) (*model.Order, error)
	UpdateStatus(req model.UpdateStatus) (*string, error)
	UpdateShipping(req model.RequestUpdateShipping) (*string, error)
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

func (h *Handler) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	UserId := mux.Vars(r)["user_id"]
	OrderId := r.URL.Query().Get("order_id")
	OrderIdUUID, err := uuid.Parse(OrderId)
	if err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	UserIDUUID, err := uuid.Parse(UserId)
	if err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get order status
	bResp, err := h.order.GetOrderStatus(UserIDUUID, OrderIdUUID)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusOK, bResp)
}

func (h *Handler) CallbackPayment(w http.ResponseWriter, r *http.Request) {
	var bReq model.RequestCallback
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Struct(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// payment success
	message, err := h.order.CallbackPayment(bReq)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusOK, message)

}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var bReq model.UpdateStatus
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Struct(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// update status
	message, err := h.order.UpdateStatus(bReq)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusOK, message)
}

func (h *Handler) UpdateShipping(w http.ResponseWriter, r *http.Request) {
	var bReq model.RequestUpdateShipping
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Struct(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	message, err := h.order.UpdateShipping(bReq)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusCreated, message)
}
