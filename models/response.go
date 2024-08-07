package models

type VariationResponse struct {
	ID       string `json:"id"`
	Size     string `json:"size" validate:"required"`
	Color    string `json:"color" validate:"required"`
	Price    int    `json:"price" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}
