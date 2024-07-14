package cart

import (
	"cart-order-service/repository/model"
	"github.com/google/uuid"
)

// cartStore is an interface that defines the methods required for managing a shopping cart.
type cartStore interface {
	// GetCartByUserID retrieves the cart for a given user.
	GetCartByUserID(bReq model.GetCartRequest) (*[]model.Cart, error)
	// GetProductDetails checks if a product exists for a given user.
	GetProductDetails(productID, userID uuid.UUID) (bool, error)
	// UpdateQty updates the quantity of a product in a user's cart.
	UpdateQty(userID, productID uuid.UUID, qty int) error
	// AddCart adds a new product to a user's cart.
	AddCart(bReq model.Cart) (*uuid.UUID, error)
	// DeleteProduct removes a product from a user's cart.
	DeleteProduct(bReq model.DeleteCartRequest) error
}

// cart is a struct that holds the store for managing a shopping cart.
type cart struct {
	store cartStore
}

// NewCart is a constructor function that returns a new cart instance.
func NewCart(store cartStore) *cart {
	return &cart{store}
}

// GetCartByUserID is a method that retrieves the cart for a given user and returns a response with the total items.
func (c *cart) GetCartByUserID(bReq model.GetCartRequest) (*[]model.Cart, error) {
	result, err := c.store.GetCartByUserID(bReq)
	if err != nil {
		return nil, err
	}

	if len(*result) == 0 {
		return nil, nil
	}

	return result, nil
}

// UpdateQty is a method that updates the quantity of a product in a user's cart or deletes the product if the quantity is 0.
func (c *cart) UpdateQty(bReq model.Cart) (string, error) {
	// if Qty is 0, delete the product from the cart
	if bReq.Qty == 0 {
		if err := c.store.DeleteProduct(model.DeleteCartRequest{
			UserID:    bReq.UserID,
			ProductID: bReq.ProductID,
		}); err != nil {
			return "", err
		}

		return "Product deleted from cart", nil
	}

	if err := c.store.UpdateQty(bReq.UserID, bReq.ProductID, bReq.Qty); err != nil {
		return "", err
	}

	return "Product updated in cart", nil
}

// AddCart is a method that adds a new product to a user's cart.
func (c *cart) AddCart(bReq model.Cart) (*uuid.UUID, error) {
	id, err := c.store.AddCart(bReq)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (c *cart) DeleteCart(bReq model.DeleteCartRequest) (string, error) {
	if err := c.store.DeleteProduct(bReq); err != nil {
		return "", err
	}

	return "Product deleted from cart", nil
}
