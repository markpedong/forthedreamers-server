package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func AddTestimonials(c *gin.Context) {
	userID := helpers.GetCurrUserToken(c).ID

	var body models.Testimonials
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var currUser models.Users
	if err := helpers.GetCurrentByID(c, &currUser, userID); err != nil {
		return
	}

	newTestimonial := models.Testimonials{
		ID:        helpers.NewUUID(),
		Title:     body.Title,
		Author:    body.Author,
		Status:    body.Status,
		ProductID: body.ProductID,
		Image:     currUser.Image,
		UserName:  currUser.Username,
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
			"created_at": v.CreatedAt,
			"image":      v.Image,
			"username":   v.UserName,
			"rating":     v.Rating,
		}
		transformedTestimonials = append(transformedTestimonials, newTestimonial)
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(transformedTestimonials))
}

func GetUserReview(c *gin.Context) {
	transformedResponse, err := helpers.GetOrderByStatus(c, true)
	if err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.JSONResponse(c, "", helpers.DataHelper(transformedResponse))
}

func AddOrderReview(c *gin.Context) {
	userID := helpers.GetCurrUserToken(c).ID

	var body models.AddTestimonials
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var cart models.CartItem
	if err := helpers.GetCurrentByID(c, &cart, body.CartID); err != nil {
		return
	}
	if cart.IsReviewed == 1 {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "You have already reviewed this product")
		return
	}

	var currUser models.Users
	if err := helpers.GetCurrentByID(c, &currUser, userID); err != nil {
		return
	}

	newTestimonial := models.Testimonials{
		ID:        helpers.NewUUID(),
		Title:     body.Description,
		Author:    currUser.FirstName + " " + currUser.LastName,
		ProductID: body.ProductID,
		Image:     currUser.Image,
		Rating:    body.Rating,
		UserName:  currUser.Username,
	}
	if err := helpers.CreateNewData(c, &newTestimonial); err != nil {
		return
	}
	cart.IsReviewed = 1
	database.DB.Save(&cart)

	helpers.JSONResponse(c, "Review added successfully")
}
