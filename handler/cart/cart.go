package cart

import (
	"cart-order-service/helper/response"
	"cart-order-service/repository/model"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// cartDto is an interface that defines the methods that our Handler struct depends on.
type cartDto interface {
	GetCartByUserID(bReq model.GetCartRequest) (*[]model.Cart, error)
	UpdateQty(bReq model.Cart) (string, error)
	AddCart(bReq model.Cart) (*uuid.UUID, error)
	DeleteCart(bReq model.DeleteCartRequest) (string, error)
}

// Handler is a struct that holds a cartDto.
type Handler struct {
	cart cartDto
}

// NewHandler is a constructor function that returns a new Handler.
func NewHandler(cart cartDto) *Handler {
	return &Handler{cart}
}

// UpdateCart is a handler function to delete and update quantity of a product from the cart.
// It first extracts the user id from the URL path, then decodes the request body into a Cart model.
// It then calls the UpdateQty method of the cartDto and sends the response back to the client.
func (h *Handler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	// Add user id from path url
	param := mux.Vars(r)
	userID := param["user_id"]

	uid, err := uuid.Parse(userID)
	if err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var bReq model.Cart
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	bReq.UserID = uid

	bResp, err := h.cart.UpdateQty(bReq)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusOK, bResp)
}

// AddCart is a handler function to add a product to the cart.
// It decodes the request body into a Cart model, checks if the quantity is greater than 0,
// then calls the AddCart method of the cartDto and sends the response back to the client.
func (h *Handler) AddCart(w http.ResponseWriter, r *http.Request) {
	var bReq model.Cart
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if bReq.Qty == 0 {
		response.HandleResponse(w, http.StatusBadRequest, "Qty must be greater than 0")
		return
	}

	bResp, err := h.cart.AddCart(bReq)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusOK, bResp)
}

// GetCartByUserID is a handler function to get a cart by user id.
// It first extracts the user id from the URL path, then decodes the request body into a GetCartRequest model.
// It then calls the GetCartByUserID method of the cartDto and sends the response back to the client.
func (h *Handler) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	// Add user id from path url
	param := mux.Vars(r)
	userID := param["user_id"]

	uid, err := uuid.Parse(userID)
	if err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var bReq model.GetCartRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var pidSlice []uuid.UUID
	pidSlice = append(pidSlice, bReq.ProductID...)

	bReq = model.GetCartRequest{
		UserID:    uid,
		ProductID: pidSlice,
	}

	bResp, err := h.cart.GetCartByUserID(bReq)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusOK, bResp)
}

func (h *Handler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	// Add user id from path url
	param := mux.Vars(r)
	userID := param["user_id"]

	uid, err := uuid.Parse(userID)
	if err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var bReq model.DeleteCartRequest
	if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
		response.HandleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	bReq.UserID = uid

	bResp, err := h.cart.DeleteCart(bReq)
	if err != nil {
		response.HandleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.HandleResponse(w, http.StatusOK, bResp)
}
