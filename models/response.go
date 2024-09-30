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
