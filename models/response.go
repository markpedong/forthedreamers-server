package models

type ProductResponse struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"Description"`
	CollectionID string              `json:"collection_id"`
	Images       []string            `json:"images"`
	Varitions    []VariationResponse `json:"variations"`
}

type VariationResponse struct {
	ID       string `json:"id"`
	Size     string `json:"size" validate:"required"`
	Color    string `json:"color" validate:"required"`
	Price    int    `json:"price" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}
