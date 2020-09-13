package middleware

import (
	"api/clients"
	"api/domains/auth"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

const bearer = "bearer"

var authorized = gin.H{
	"error": true,
	"message": "Unauthorized",
	"status_code": http.StatusUnauthorized,
}

func ValidateAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, authorized)
			c.Abort()
			return
		}

		if len(authHeader) < len(bearer) {
			c.JSON(http.StatusUnauthorized, authorized)
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader[len(bearer):], " ", "", -1)

		// Check is token is expired itself
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, authorized)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user := auth.User{
				ID:    0,
				Email: "",
			}
			err := clients.ClientOrm.Table("users").
				Select("id","email").
				Where("id", "=", claims["id"]).
				First(&user)

			if err != nil {
				c.JSON(http.StatusUnauthorized, authorized)
				c.Abort()
				return
			}

			c.Set("user", user)
			c.Next()
			return
		} else {
			c.JSON(http.StatusUnauthorized, authorized)
			c.Abort()
			return
		}
	}
}