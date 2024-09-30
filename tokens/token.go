package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(bytes)
}

func CreateAndSignJWT(userID *string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":      userID,
		"randomValue": GenerateRandomString(5),
		"ttl":         time.Now().Add(time.Hour * 12).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("HMAC_SECRET")))

}
