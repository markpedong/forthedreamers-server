package models

type VariationResponse struct {
	ID       string `json:"id"`
	Size     string `json:"size" validate:"required"`
	Color    string `json:"color" validate:"required"`
	Price    int    `json:"price" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

type UsersResponse struct {
	ID        string `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Image     string `json:"image" validate:"required"`
	Username  string `json:"username"`
}

type CredentialResponse struct {
	UserInfo UsersResponse `json:"userInfo"`
	Token    string        `json:"token"`
}

type AddressItemReponse struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Address   string `json:"address" validate:"required"`
	IsDefault int    `json:"is_default" gorm:"default:0"`
}

type CartItemResponse struct {
	ID          string   `json:"id" gorm:"primaryKey"`
	Quantity    int      `json:"quantity" validate:"required"`
	ProductName string   `json:"name" validate:"required"`
	ProductID   string   `json:"product_id" validate:"required"`
	Size        string   `json:"size,omitempty" validate:"required"`
	Color       string   `json:"color,omitempty" validate:"required"`
	Price       int      `json:"price" validate:"required"`
	Image       []string `json:"images" validate:"required"`
}

type ItemResponse struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Quantity    int    `json:"quantity" validate:"required"`
	ProductName string `json:"name" validate:"required"`
	ProductID   string `json:"product_id" validate:"required"`
	Size        string `json:"size,omitempty" validate:"required"`
	Color       string `json:"color,omitempty" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	Image       string `json:"images" validate:"required"`
	CreatedAt   int    `json:"created_at"`
	IsReviewed  int    `json:"is_reviewed"`
}

type OrderItemResponse struct {
	ID            string             `json:"id"`
	TotalPrice    int                `json:"total_price" validate:"required"`
	PaymentMethod int                `json:"payment_method" validate:"required"`
	Address       AddressItemReponse `json:"address" validate:"required"`
	Items         []ItemResponse     `json:"items"`
	CreatedAt     int                `json:"created_at"`
	Status        int                `json:"status"`
}
