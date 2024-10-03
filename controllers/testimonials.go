package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func AddTestimonials(c *gin.Context) {
	var body models.Testimonials
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	newTestimonial := models.Testimonials{
		ID:        helpers.NewUUID(),
		Title:     body.Title,
		Author:    body.Author,
		Status:    body.Status,
		ProductID: body.ProductID,
	}
	if err := helpers.CreateNewData(c, &newTestimonial); err != nil {
		return
	}

	helpers.JSONResponse(c, "")
}

func GetTestimonials(c *gin.Context) {
	var testimonials []models.Testimonials
	if err := database.DB.Find(&testimonials).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(testimonials))
}

func PublicTestimonials(c *gin.Context) {
	var testimonials []models.Testimonials
	if err := database.DB.Where("status = ?", 1).Find(&testimonials).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	transformedTestimonials := []map[string]interface{}{}
	for _, v := range testimonials {
		newTestimonial := map[string]interface{}{
			"id":         v.ID,
			"product_id": v.ProductID,
			"author":     v.Author,
			"title":      v.Title,
		}
		transformedTestimonials = append(transformedTestimonials, newTestimonial)
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(transformedTestimonials))
}
