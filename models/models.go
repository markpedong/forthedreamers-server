package models

import (
	// "github.com/lib/pq"
	"github.com/lib/pq"
	"gorm.io/plugin/soft_delete"
)

type Users struct {
	ID        string                `json:"id" gorm:"primaryKey"`
	FirstName string                `json:"first_name" validate:"required"`
	LastName  string                `json:"last_name" validate:"required"`
	Phone     string                `json:"phone" validate:"required"`
	Address   string                `json:"address" validate:"required"`
	Email     string                `json:"email" validate:"required"`
	Image     string                `json:"image" validate:"required"`
	Username  string                `json:"username"`
	Password  string                `json:"password"`
	CreatedAt int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
}

type Collection struct {
	ID        string                `json:"id" gorm:"primaryKey"`
	Name      string                `json:"name" validate:"required"`
	Images    pq.StringArray        `json:"images" gorm:"type:text[]"  validate:"required"`
	Products  []Product             `json:"products" gorm:"foreignKey:CollectionID"`
	CreatedAt int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
}

type Product struct {
	ID           string                `json:"id" gorm:"primaryKey"`
	Name         string                `json:"name" validate:"required"`
	Description  string                `json:"Description" validate:"required"`
	CollectionID string                `json:"collection_id"`
	Price        int                   `json:"price" validate:"required"`
	Quantity     int                   `json:"quantity" validate:"required"`
	Size         string                `json:"size" validate:"required"`
	Color        string                `json:"color" validate:"required"`
	CreatedAt    int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    soft_delete.DeletedAt `json:"-"`
}
