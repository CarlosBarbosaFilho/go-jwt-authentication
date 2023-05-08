package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/CarlosBarbosaFilho/go-jwt-authentication/controllers/response"
	"github.com/CarlosBarbosaFilho/go-jwt-authentication/initializers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAtuthentication(ctx *gin.Context) {
	fmt.Println("In Middleware")
	//Get cookie of request
	tokenString, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return ([]byte(os.Getenv("SECRET_ACCESS"))), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		//Find user with user sub
		var user response.DataUser
		//initializers.DB.First(&user, claims["sub"])
		initializers.DB.Raw("SELECT * FROM users WHERE id = ?", claims["sub"]).Scan(&user)
		if user.ID == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		//Attach to request
		ctx.Set("user", user)

		//Continue
		ctx.Next()
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

}
