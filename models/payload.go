package models

import (
	"github.com/lib/pq"
)

type CollectionPayload struct {
	ID     string         `json:"id"`
	Name   string         `json:"name" validate:"required"`
	Images pq.StringArray `json:"images" gorm:"type:text[]"  validate:"required"`
}

type UserPayload struct {
	ID        string        `json:"id"`
	Image     string        `json:"image"`
	FirstName string        `json:"first_name" validate:"required"`
	LastName  string        `json:"last_name" validate:"required"`
	Phone     string        `json:"phone" validate:"required"`
	Address   []AddressItem `json:"address"`
	Email     string        `json:"email" validate:"required"`
	Username  string        `json:"username" validate:"required"`
	Password  string        `json:"password" validate:"required"`
}

type UpdateUserPayload struct {
	ID          string `json:"id"`
	Image       string `json:"image"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Username    string `json:"username" validate:"required"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type ProductVariationPayload struct {
	ID        string `json:"id"`
	Size      string `json:"size" validate:"required"`
	Color     string `json:"color" validate:"required"`
	Price     int    `json:"price" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required"`
	ProductID string `json:"product_id"`
}

type ProductPayload struct {
	ID           string         `json:"id"`
	Name         string         `json:"name" validate:"required"`
	Description  string         `json:"Description" validate:"required"`
	CollectionID string         `json:"collection_id"`
	Images       pq.StringArray `json:"images" gorm:"type:text[]"  validate:"required"`
	Features     pq.StringArray `json:"features" gorm:"type:text[]"  validate:"required"`
}

type AddTestimonials struct {
	CartID      string `json:"cart_id" validate:"required"`
	ProductID   string `json:"product_id" validate:"required"`
	Rating      int    `json:"rating" validate:"required"`
	Description string `json:"description"`
}
