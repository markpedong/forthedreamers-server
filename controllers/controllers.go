package controllers

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/forthedreamers-server/cloudinary"
	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/forthedreamers-server/tokens"
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

func Login(ctx *gin.Context) {
	var body struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	if err := helpers.BindValidateJSON(ctx, &body); err != nil {
		return
	}

	var existingUser models.Users
	if err := database.DB.First(&existingUser, "username = ?", body.UserName).Error; err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	notValid, msg := VerifyPassword(existingUser.Password, body.Password)
	if notValid {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, msg)
		return
	}

	token, err := tokens.CreateAndSignJWT(&existingUser)
	if err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	tokens.SetCookie(ctx, token)
	userRes := map[string]interface{}{
		"token":    token,
		"userInfo": existingUser,
	}

	ctx.Redirect(http.StatusFound, "/app/users")
	helpers.JSONResponse(ctx, "", helpers.DataHelper(userRes))
}

func UploadImage(ctx *gin.Context) {
	form, err := ctx.FormFile("file")

	if err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	uploadResult, err := cloudinary.CloudinaryService.Upload.Upload(ctx, form, uploader.UploadParams{Folder: "forthedreamers"})
	if err != nil {
		helpers.ErrJSONResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	imageRes := map[string]interface{}{
		"url":      uploadResult.URL,
		"fileName": uploadResult.OriginalFilename,
		"size":     uploadResult.Bytes,
	}

	helpers.JSONResponse(ctx, "upload successful!", helpers.DataHelper(imageRes))
}
