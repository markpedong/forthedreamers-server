package controllers

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/forthedreamers-server/cloudinary"
	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
)

func VerifyPassword(expectedHashedPassword, givenPassword string) (bool, string) {
	// err := bcrypt.CompareHashAndPassword([]byte(expectedHashedPassword), []byte(givenPassword))
	err := expectedHashedPassword == givenPassword

	switch {
	case err:
		return false, "Password matched!"
	// case errors.Is(_, bcrypt.ErrMismatchedHashAndPassword):
	// 	return false, "Password is incorrect!"
	case !err:
		return true, "Password is incorrect!"
	default:
		// fmt.Printf("Password verification error: %s\n", err)
		return true, "Failed to verify password"
	}
}

func UploadImage(c *gin.Context) {
	form, err := c.FormFile("file")

	if err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	uploadResult, err := cloudinary.CloudinaryService.Upload.Upload(c, form, uploader.UploadParams{
		Folder:         "forthedreamers",
		Transformation: "f_webp,q_auto:good,fl_lossy,c_fit",
	})

	if err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	imageRes := map[string]interface{}{
		"url":      uploadResult.URL,
		"fileName": uploadResult.OriginalFilename,
		"size":     uploadResult.Bytes,
	}

	helpers.JSONResponse(c, "upload successful!", helpers.DataHelper(imageRes))
}

func SignUp(c *gin.Context) {
	var body struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required"`
		Password  string `json:"password" validate:"required"`
		Username  string `json:"username" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	newUser := models.Users{
		ID:        helpers.NewUUID(),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
		Username:  body.Username,
	}
	existingUser := models.Users{}
	if err := database.DB.Where("email = ? OR username = ?", body.Email, body.Username).First(&existingUser).Error; err == nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, "email or username already exists")
		return
	}
	if err := helpers.CreateNewData(c, &newUser); err != nil {
		return
	}

	userRes := helpers.UserGetTokenResponse(c, &newUser)
	helpers.JSONResponse(c, "", helpers.DataHelper(userRes))
}
