package tokens

import (
	"os"
	"time"

	"github.com/forthedreamers-server/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateAndSignJWT(user *models.Users) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"ttl":    time.Now().Add(time.Hour * 24 * 100).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("HMAC_SECRET")))

}

func SetCookie(ctx *gin.Context, token string) {
	// server side only
	// ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("Auth", token, 3600*24, "", "", true, false)
	// http.SetCookie(ctx.Writer, &http.Cookie{
	// 	Name:     "Auth",
	// 	Value:    token,
	// 	Path:     "/",
	// 	Domain:   "",
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteNoneMode,
	// })
}
