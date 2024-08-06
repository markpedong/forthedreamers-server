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
	Status    int                   `json:"status" gorm:"default:0"`
	Role      string                `json:"role" gorm:"default:USER"`
}

type Collection struct {
	ID        string                `json:"id" gorm:"primaryKey"`
	Name      string                `json:"name" validate:"required"`
	Images    pq.StringArray        `json:"images" gorm:"type:text[]" validate:"required"`
	Products  []Product             `json:"products" gorm:"foreignKey:CollectionID"`
	CreatedAt int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
}

type ProductVariation struct {
	ID        string                `json:"id" gorm:"primaryKey"`
	Size      string                `json:"size" validate:"required"`
	Color     string                `json:"color" validate:"required"`
	Price     int                   `json:"price" validate:"required"`
	Quantity  int                   `json:"quantity" validate:"required"`
	CreatedAt int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
	ProductID string                `json:"product_id"`
	Status    int                   `json:"status" gorm:"default:0"`
}

type Product struct {
	ID           string                `json:"id" gorm:"primaryKey"`
	Name         string                `json:"name" validate:"required"`
	Description  string                `json:"description" validate:"required"`
	CollectionID string                `json:"collection_id"`
	Variations   []ProductVariation    `json:"variations" gorm:"foreignKey:ProductID"`
	Images       pq.StringArray        `json:"images" gorm:"type:text[]"`
	Features     pq.StringArray        `json:"features" gorm:"type:text[]"`
	CreatedAt    int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    soft_delete.DeletedAt `json:"-"`
	Status       int                   `json:"status" gorm:"default:0"`
}

type WebsiteData struct {
	ID            string                `json:"id" gorm:"primaryKey"`
	WebsiteName   string                `json:"website_name" validate:"required"`
	PromoText     string                `json:"promo_text" validate:"required"`
	MarqueeText   string                `json:"marquee_text" validate:"required"`
	NewsText      string                `json:"news_text" validate:"required"`
	LandingImage1 string                `json:"landing_image1" validate:"required"`
	LandingImage2 string                `json:"landing_image2" validate:"required"`
	LandingImage3 string                `json:"landing_image3" validate:"required"`
	CreatedAt     int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     soft_delete.DeletedAt `json:"-"`
}
