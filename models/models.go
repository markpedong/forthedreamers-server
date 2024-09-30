package models

import (
	// "github.com/lib/pq"
	"github.com/lib/pq"
	"gorm.io/plugin/soft_delete"
)

type AddressItem struct {
	ID        string `json:"id" gorm:"primaryKey"`
	IsDefault int    `json:"is_default" gorm:"default:0"`
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Address   string `json:"address" validate:"required"`
}

type Users struct {
	ID        string                `json:"id" gorm:"primaryKey"`
	FirstName string                `json:"first_name" validate:"required"`
	LastName  string                `json:"last_name" validate:"required"`
	Phone     string                `json:"phone"`
	Address   []AddressItem         `json:"address" gorm:"foreignKey:UserID"`
	Email     string                `json:"email" validate:"required"`
	Image     string                `json:"image"`
	Username  string                `json:"username" validate:"required"`
	Password  string                `json:"password" validate:"required"`
	CreatedAt int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
	Status    int                   `json:"status" gorm:"default:0"`
	Token     string                `json:"token"`
	CartItems []UserCart            `json:"cart_items" gorm:"foreignKey:UserID"`
}

type Collection struct {
	ID        string                `json:"id" gorm:"primaryKey"`
	Name      string                `json:"name" validate:"required"`
	Images    pq.StringArray        `json:"images" gorm:"type:text[]" validate:"required"`
	Products  []Product             `json:"products" gorm:"foreignKey:CollectionID"`
	CreatedAt int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
	Status    int                   `json:"status" gorm:"default:0"`
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
	Testimonials []Testimonials        `json:"testimonials" gorm:"foreignKey:ProductID"`
	CartItems    []CartItem            `json:"cart_items" gorm:"foreignKey:ProductID"`
	Images       pq.StringArray        `json:"images" gorm:"type:text[]"`
	Features     pq.StringArray        `json:"features" gorm:"type:text[]"`
	CreatedAt    int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    soft_delete.DeletedAt `json:"-"`
	Status       int                   `json:"status" gorm:"default:0"`
}

type WebsiteData struct {
	ID              string                `json:"id" gorm:"primaryKey"`
	WebsiteName     string                `json:"website_name" validate:"required"`
	PromoText       string                `json:"promo_text" validate:"required"`
	MarqueeText     string                `json:"marquee_text" validate:"required"`
	NewsText        string                `json:"news_text" validate:"required"`
	LandingImage1   string                `json:"landing_image1" validate:"required"`
	LandingImage2   string                `json:"landing_image2" validate:"required"`
	LandingImage3   string                `json:"landing_image3" validate:"required"`
	DefaultPageSize int                   `json:"default_pageSize" validate:"required"`
	CreatedAt       int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       soft_delete.DeletedAt `json:"-"`
}

type Testimonials struct {
	ID        string                `json:"id" gorm:"primaryKey"`
	Title     string                `json:"title" validate:"required"`
	Author    string                `json:"author" validate:"required"`
	Status    int                   `json:"status" gorm:"default:0"`
	ProductID string                `json:"product_id" validate:"required"`
	CreatedAt int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `json:"-"`
}

type CartItem struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Quantity    int    `json:"quantity" validate:"required"`
	ProductID   string `json:"product_id" validate:"required"`
	VariationID string `json:"variation_id" validate:"required"`
	OrderItemID string `json:"-" validate:"required"`
	Status      int    `json:"status" gorm:"default:0"` // 0 - cart, 1 - ordered, 2 - delivered, 3 - canceled/returned
}

type UserCart struct {
	UserID     string `gorm:"primaryKey"`
	CartItemID string `gorm:"primaryKey"`
}

type OrderItem struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	AddressID string     `json:"address_id" validate:"required"`
	Items     []CartItem `json:"items" validate:"required" gorm:"foreignKey:OrderItemID"`
	Price     int        `json:"price" validate:"required"`
	Status    int        `json:"status" gorm:"default:0"` // 0 - pending, 1 - in transit, 2 - out for delivery, 3 - delivered
}
