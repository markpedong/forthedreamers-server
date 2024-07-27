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
	ID        string `json:"id"`
	Image     string `json:"image" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Address   string `json:"address" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type ProductPayload struct {
	ID           string         `json:"id"`
	Name         string         `json:"name" validate:"required"`
	Description  string         `json:"Description" validate:"required"`
	CollectionID string         `json:"collection_id"`
	Price        int            `json:"price" validate:"required"`
	Quantity     int            `json:"quantity" validate:"required"`
	Images       pq.StringArray `json:"images" gorm:"type:text[]"  validate:"required"`
	Size         pq.StringArray `json:"sizes" gorm:"type:text[]" validate:"required"`
	Color        pq.StringArray `json:"colors" gorm:"type:text[]" validate:"required"`
}
