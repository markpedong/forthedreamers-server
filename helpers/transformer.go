package helpers

import (
	"net/http"
	"strings"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/models"
	"github.com/forthedreamers-server/tokens"
	"github.com/gin-gonic/gin"
)

func UserGetTokenResponse(c *gin.Context, user *models.Users) models.CredentialResponse {
	token, err := tokens.CreateAndSignJWT(&user.ID)
	if err != nil {
		ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return models.CredentialResponse{}
	}
	user.Token = token
	database.DB.Save(&user)

	newUserResponse := models.UsersResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Image:     user.Image,
		Username:  user.Username,
		Phone:     user.Phone,
	}

	partToken := strings.Split(token, ".")
	userRes := models.CredentialResponse{
		UserInfo: newUserResponse,
		Token:    partToken[1],
	}

	return userRes
}
