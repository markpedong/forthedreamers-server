package tokens

import (
	"net/http"
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
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Auth", token, 3600*24*100, "", "", false, true)
}
