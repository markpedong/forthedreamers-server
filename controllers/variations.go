package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func PublicVariations(ctx *gin.Context) {
	var body struct {
		ID string `json:"product_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currVariations []models.ProductVariation
	if err := database.DB.Where("status = ?", 1).Order("created_at DESC").Find(&currVariations, "product_id = ?", body.ID).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
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

	helpers.JSONResponse(ctx, "", helpers.DataHelper(productVariations))
}

func GetVariations(ctx *gin.Context) {
	var body struct {
		ProductID string `json:"product_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var variations []models.ProductVariation
	if err := database.DB.Find(&variations, "product_id = ?", body.ProductID).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(variations))
}

func UpdateVariations(ctx *gin.Context) {
	var body models.ProductVariationPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var currVariation models.ProductVariation
	if err := helpers.GetCurrentByID(ctx, &currVariation, body.ID); err != nil {
		return
	}

	helpers.UpdateByModel(ctx, &currVariation, models.ProductVariation{Size: body.Size, Color: body.Color, Price: body.Price, Quantity: body.Quantity})
	helpers.JSONResponse(ctx, "")
}

func ToggleVariations(ctx *gin.Context) {
	var body struct {
		ID string `json:"ID" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := helpers.ToggleModelByID(ctx, &models.ProductVariation{}, body.ID); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}

func DeleteVariations(ctx *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currVariation models.ProductVariation
	if err := helpers.GetCurrentByID(ctx, &currVariation, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE VARIATION IS EXISTENT
	helpers.DeleteByModel(ctx, &currVariation)
	helpers.JSONResponse(ctx, "", helpers.DataHelper(&currVariation))
}

func AddVariations(ctx *gin.Context) {
	var body models.ProductVariationPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currProduct models.Product
	if err := helpers.GetCurrentByID(ctx, &currProduct, body.ProductID); err != nil {
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
	if err := helpers.CreateNewData(ctx, &newVariation); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}
