package controllers

import (
	"net/http"
	"strconv"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PublicCollections(c *gin.Context) {
	var body struct {
		PageSize string `form:"page_size,omitempty"`
		Page     string `form:"page,omitempty"`
	}
	if err := helpers.BindValidateQuery(c, &body); err != nil {
		return
	}

	var collections []models.Collection
	db := database.DB.Where("status = ?", 1)

	if body.PageSize != "" && body.Page != "" {
		pageSize, err1 := strconv.Atoi(body.PageSize)
		page, err2 := strconv.Atoi(body.Page)

		if err1 == nil && err2 == nil {
			db = db.Limit(pageSize).Offset((page - 1) * pageSize)
		}
	}

	if err := db.Find(&collections).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	transformedCollections := make([]map[string]interface{}, len(collections))
	for i, v := range collections {
		transformedCollections[i] = map[string]interface{}{
			"id":     v.ID,
			"name":   v.Name,
			"images": []string{v.Images[0]},
		}
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(transformedCollections))
}

func AddCollection(c *gin.Context) {
	var body models.CollectionPayload
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	if len(body.Images) < 1 {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "images is required")
		return
	}

	newCollection := models.Collection{
		ID:     helpers.NewUUID(),
		Name:   body.Name,
		Images: body.Images,
	}
	if err := helpers.CreateNewData(c, &newCollection); err != nil {
		return
	}

	helpers.JSONResponse(c, "")
}

func GetCollection(c *gin.Context) {
	var collections []models.Collection

	// NO NEED TO HANDLE ERROR HERE BECAUSE COLLECTION IS EXISTENT
	helpers.GetTableByModel(c, &collections)
	helpers.JSONResponse(c, "", helpers.DataHelper(&collections))
}

func UpdateCollection(c *gin.Context) {
	var body models.CollectionPayload
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currCollection models.Collection
	if err := helpers.GetCurrentByID(c, &currCollection, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE COLLECTION IS EXISTENT
	helpers.UpdateByModel(c, &currCollection, models.Collection{Name: body.Name, Images: body.Images})
	helpers.JSONResponse(c, "")

}

func DeleteCollection(c *gin.Context) {
	var body struct {
		ID string `json:"id" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currCollection models.Collection
	if err := helpers.GetCurrentByID(c, &currCollection, body.ID); err != nil {
		return
	}

	// NO NEED TO HANDLE ERROR HERE BECAUSE COLLECTION IS EXISTENT
	helpers.DeleteByModel(c, &currCollection)
	helpers.JSONResponse(c, "")
}

func ToggleCollections(c *gin.Context) {
	var body struct {
		ID string `json:"ID" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := helpers.ToggleModelByID(c, &models.Collection{}, body.ID); err != nil {
		return
	}

	helpers.JSONResponse(c, "")
}

func GetCollectionByID(c *gin.Context) {
	var body struct {
		ID string `form:"id" validate:"required"`
	}
	if err := helpers.BindValidateQuery(c, &body); err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var collection models.Collection
	if err := database.DB.Where("status = ?", 1).First(&collection, "id = ?", body.ID).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var collectionsProducts []models.Product
	if err := database.DB.
		Preload("Variations", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("status = ?", 1).
				Order("created_at DESC")
		}).
		Find(&collectionsProducts, "collection_id = ?", collection.ID).
		Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var filteredProducts []map[string]interface{}
	for _, product := range collectionsProducts {
		var filteredVariations []map[string]interface{}
		for _, w := range product.Variations {
			variation := map[string]interface{}{
				"id":    w.ID,
				"size":  w.Size,
				"color": w.Color,
				"price": w.Price,
			}
			filteredVariations = append(filteredVariations, variation)
		}

		filteredProduct := map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"images":      product.Images,
			"variations":  filteredVariations,
		}
		filteredProducts = append(filteredProducts, filteredProduct)
	}

	newCollection := map[string]interface{}{
		"id":       collection.ID,
		"name":     collection.Name,
		"images":   collection.Images,
		"products": filteredProducts,
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(newCollection))
}
