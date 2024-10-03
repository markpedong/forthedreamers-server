package controllers

import (
	"net/http"

	"github.com/forthedreamers-server/database"
	"github.com/forthedreamers-server/helpers"
	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not complete authentication"})
		return
	}

	oauthToken := user.AccessToken

	script := `
        <html>
        <head>
            <script>
                window.opener.location.href = "http://localhost:6600/login?otp=` + oauthToken + `";
                window.close();
            </script>
        </head>
        <body></body>
        </html>
    `
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
		helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
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
