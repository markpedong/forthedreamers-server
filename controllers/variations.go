package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func PublicVariations(c *gin.Context) {
	var body struct {
		ID string `json:"product_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currVariations []models.ProductVariation
	if err := database.DB.Where("status = ?", 1).Order("created_at DESC").Find(&currVariations, "product_id = ?", body.ID).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var productVariations []models.VariationResponse
	for _, v := range currVariations {
		newProductVariation := models.VariationResponse{
			ID:       v.ID,
			Size:     v.Size,
			Color:    v.Color,
			Price:    v.Price,
			Quantity: v.Quantity,
		}

		productVariations = append(productVariations, newProductVariation)
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(productVariations))
}

func GetVariations(c *gin.Context) {
	var body struct {
		ProductID string `json:"product_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var variations []models.ProductVariation
	if err := database.DB.Find(&variations, "product_id = ?", body.ProductID).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(variations))
}

func UpdateVariations(c *gin.Context) {
	var body models.ProductVariationPayload
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var currVariation models.ProductVariation
	if err := helpers.GetCurrentByID(c, &currVariation, body.ID); err != nil {
		return
	}

	helpers.UpdateByModel(c, &currVariation, models.ProductVariation{Size: body.Size, Color: body.Color, Price: body.Price, Quantity: body.Quantity})
	helpers.JSONResponse(c, "")
}

func ToggleVariations(c *gin.Context) {
	var body struct {
		ID string `json:"ID" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := helpers.ToggleModelByID(c, &models.ProductVariation{}, body.ID); err != nil {
		return
	}

	helpers.JSONResponse(c, "")
}

func DeleteVariations(c *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currVariation models.ProductVariation
	if err := helpers.GetCurrentByID(c, &currVariation, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE VARIATION IS EXISTENT
	helpers.DeleteByModel(c, &currVariation)
	helpers.JSONResponse(c, "", helpers.DataHelper(&currVariation))
}

func AddVariations(c *gin.Context) {
	var body models.ProductVariationPayload
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currProduct models.Product
	if err := helpers.GetCurrentByID(c, &currProduct, body.ProductID); err != nil {
		return
	}

	newVariation := models.ProductVariation{
		ID:        helpers.NewUUID(),
		Size:      body.Size,
		Color:     body.Color,
		Price:     body.Price,
		Quantity:  body.Quantity,
		ProductID: currProduct.ID,
	}
	if err := helpers.CreateNewData(c, &newVariation); err != nil {
		return
	}

	helpers.JSONResponse(c, "")
}
