package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/forthedreamers-server/tokens"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

func GoogleLogin(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func GoogleCallback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, "Error: %s", err.Error())
		return
	}

	var existingUser models.Users
	err = database.DB.First(&existingUser, "email = ?", user.Email).Error

	if err == gorm.ErrRecordNotFound {
		existingUser = models.Users{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Image:     user.AvatarURL,
			Username:  strings.ReplaceAll(user.NickName, " ", ""),
			ID:        user.UserID,
		}
	} else if err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := tokens.CreateAndSignJWT(&user.UserID)
	if err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	existingUser.Token = token
	if err := database.DB.Save(&existingUser).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	script := fmt.Sprintf(`
		<script>
			if (window.opener) {
				window.opener.postMessage({ data: { userInfo: %s, token: '%s', code: 200, message: 'Logged in successfully' }}, '*');
			}
			window.close();
		</script>
	`, helpers.ToJSON(existingUser), strings.Split(token, ".")[1])

	c.Data(http.StatusOK, "text/html", []byte(script))
}

func Login(c *gin.Context) {
	var body struct {
		UserName string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	var existingUser models.Users
	if err := database.DB.First(&existingUser, "username = ?", body.UserName).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "user doesn't exist")
		return
	}

	notValid, msg := VerifyPassword(existingUser.Password, body.Password)
	if notValid {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, msg)
		return
	}

	userRes := helpers.UserGetTokenResponse(c, &existingUser)
	helpers.JSONResponse(c, "Logged in successfully!", helpers.DataHelper(userRes))
}
