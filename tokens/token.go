package tokens

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateAndSignJWT(userID *string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 12).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    expirationTime,
	})

	return token.SignedString([]byte(os.Getenv("HMAC_SECRET")))

}
