package middlewares

import (
	"../models"
	"../settings"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func IsAuthorized() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token, _ := c.Cookie("token")
		if header == "" {
			header = fmt.Sprintf("cookie %s", token)
		}
		headerSplit := strings.Split(header, " ")
		if len(headerSplit) > 1 && len(headerSplit[1]) > 10 {
			tokenString := headerSplit[1]

			// Parse takes the token string and a function for looking up the key. The latter is especially
			// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
			// head of the token to identify which key to use, but the parsed token (head and claims) is provided
			// to the callback, providing flexibility.
			token, err := jwt.ParseWithClaims(tokenString, &models.AuthUser{}, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return settings.JwtKey, nil
			})
			if claims, ok := token.Claims.(*models.AuthUser); ok && token.Valid {
				claims.Token = tokenString
				c.Set("user_claims", claims)
				c.Set("user", claims)
			} else {
				fmt.Println(err)
			}
		}
		c.Next()
	})
}
