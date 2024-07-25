package models

import (
	// "github.com/lib/pq"
	"gorm.io/plugin/soft_delete"
)

type Users struct {
	ID          string                `json:"id" gorm:"primaryKey"`
	FirstName   string                `json:"first_name" validate:"required"`
	LastName    string                `json:"last_name" validate:"required"`
	Phone       string                `json:"phone" validate:"required"`
	Address     string                `json:"address" validate:"required"`
	Description string                `json:"description" validate:"required"`
	Email       string                `json:"email" validate:"required"`
	Username    string                `json:"username"`
	Password    string                `json:"password"`
	CreatedAt   int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   soft_delete.DeletedAt `json:"-"`
}
