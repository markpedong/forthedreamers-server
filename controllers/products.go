package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateNewProduct(payload *models.ProductPayload) models.Product {
	return models.Product{
		ID:           helpers.NewUUID(),
		Name:         payload.Name,
		Description:  payload.Description,
		CollectionID: payload.CollectionID,
		Features:     payload.Features,
	}
}

func PublicProductDetails(ctx *gin.Context) {
	var body struct {
		ID string `json:"product_id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currProduct models.Product
	if err := helpers.GetCurrentByID(ctx, &currProduct, body.ID); err != nil {
		return
	}

	filteredProduct := map[string]interface{}{
		"id":          currProduct.ID,
		"name":        currProduct.Name,
		"description": currProduct.Description,
		"images":      currProduct.Images,
		"features":    currProduct.Features,
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(filteredProduct))
}

func PublicProducts(ctx *gin.Context) {
	var body struct {
		Search   string `form:"search"`
		PageSize int    `form:"page_size"`
		Page     int    `form:"page"`
	}
	if err := helpers.BindValidateQuery(ctx, &body); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	query := database.DB.Preload("Variations", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", 1).Order("created_at DESC")
	}).Where("status = ?", 1)

	if body.PageSize > 0 && body.Page > 0 {
		query.Offset((body.Page - 1) * body.PageSize).Limit(body.PageSize)
	}

	if body.Search != "" {
		query = query.Where("name ILIKE ?", "%"+body.Search+"%")
	}

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	filteredProducts := make([]map[string]interface{}, len(products))
	for i, product := range products {
		filteredVariations := make([]map[string]interface{}, len(product.Variations))
		for j, w := range product.Variations {
			filteredVariations[j] = map[string]interface{}{
				"id":    w.ID,
				"size":  w.Size,
				"color": w.Color,
				"price": w.Price,
			}
		}

		filteredProducts[i] = map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"images":      product.Images,
			"variations":  filteredVariations,
		}
	}

	helpers.JSONResponse(ctx, "", helpers.DataHelper(filteredProducts))
}

func GetProducts(ctx *gin.Context) {
	var products []models.Product

	//NO NEED TO HANDLE ERROR BECAUSE PRODUCT IS EXISTENT
	helpers.GetTableByModel(ctx, &products)
	helpers.JSONResponse(ctx, "", helpers.DataHelper(products))
}

func AddProducts(ctx *gin.Context) {
	var body models.ProductPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	newProduct := models.Product{
		ID:           helpers.NewUUID(),
		Name:         body.Name,
		Images:       body.Images,
		Description:  body.Description,
		CollectionID: body.CollectionID,
		Features:     body.Features,
	}

	if err := helpers.CreateNewData(ctx, &newProduct); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}

func UpdateProducts(ctx *gin.Context) {
	var body models.ProductPayload
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var currProduct models.Product
	if err := helpers.GetCurrentByID(ctx, &currProduct, body.ID); err != nil {
		return
	}

	//NO NEED TO HANDLE ERROR BECAUSE PRODUCT IS EXISTENT
	helpers.UpdateByModel(ctx, &currProduct, models.Product{
		Name:         body.Name,
		Images:       body.Images,
		Description:  body.Description,
		CollectionID: body.CollectionID,
		Features:     body.Features,
	})
	helpers.JSONResponse(ctx, "")
}

func DeleteProducts(ctx *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var currProduct models.Product
	if err := helpers.GetCurrentByID(ctx, &currProduct, body.ID); err != nil {
		return
	}

	var variations []models.ProductVariation
	if err := database.DB.Find(&variations, "product_id = ?", currProduct.ID).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range variations {
		helpers.DeleteByModel(ctx, &v)
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE CURRPRODUCT IS EXISTENT
	helpers.DeleteByModel(ctx, &currProduct)
	helpers.JSONResponse(ctx, "")
}

func ToggleProducts(ctx *gin.Context) {
	var body struct {
		ID string `json:"ID" validate:"required"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := helpers.ToggleModelByID(ctx, &models.Product{}, body.ID); err != nil {
		return
	}

	helpers.JSONResponse(ctx, "")
}
