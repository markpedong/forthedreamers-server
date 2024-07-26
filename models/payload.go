package models

import (
	"github.com/lib/pq"
)

type CollectionPayload struct {
	Name   string         `json:"name" validate:"required"`
	Images pq.StringArray `json:"images" gorm:"type:text[]"  validate:"required"`
}

type UserPayload struct {
	Image     string `json:"image" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Address   string `json:"address" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}
