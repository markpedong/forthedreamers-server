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

func RequestEmailOTP(c *gin.Context) {
	var body struct {
		Email string `json:"email" validate:"required"`
	}
	if err := helpers.BindValidateJSON(c, &body); err != nil {
		return
	}

	if err := database.DB.Select("email").First(&models.Users{}, "email = ?", body.Email).Error; err != nil {
		helpers.ErrJSONResponse(c, http.StatusBadRequest, "email doesn't exist")
		return
	}

	// if err := sendMailSimple(body.Email); err != nil {
	// 	helpers.ErrJSONResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	helpers.JSONResponse(c, "Successfully sent!")
}

// func sendMailSimple(mail string) error {
// 	auth := smtp.PlainAuth(
// 		"",
// 		"forthedreamersforthedreamers@gmail.com",
// 		os.Getenv("APP_PASSWORD"),
// 		"smtp.gmail.com",
// 	)

// 	msg := []byte("To: " + mail + "\r\n" +
// 		"From: forthedreamersforthedreamers@gmail.com\r\n" +
// 		"Subject: Reset Password OTP Email ( DO NOT SHARE WITH ANYONE )\r\n" +
// 		"\r\n" +
// 		"This is a test mail")

// 	if err := smtp.SendMail(
// 		"smtp.gmail.com:587",
// 		auth,
// 		"forthedreamersforthedreamers@gmail.com",
// 		[]string{
// 			mail,
// 		},
// 		msg,
// 	); err != nil {
// 		return err
// 	}

// 	return nil
// }

func VerifyOTP(c *gin.Context) {}

func SetNewPassword(c *gin.Context) {}

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
	helpers.JSONResponse(c, "User created successfully", helpers.DataHelper(userRes))
}
