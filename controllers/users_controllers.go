package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/CarlosBarbosaFilho/go-jwt-authentication/controllers/request"
	"github.com/CarlosBarbosaFilho/go-jwt-authentication/controllers/response"
	"github.com/CarlosBarbosaFilho/go-jwt-authentication/initializers"
	"github.com/CarlosBarbosaFilho/go-jwt-authentication/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(ctx *gin.Context) {

	// Get the email/password of request body
	request := request.UserRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Erro to read request body"})
		return
	}
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Create user on database with your credentials
	user := models.User{Name: request.Name, Email: request.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Erro to create user"})
		return
	}
	// Response the reques received
	var response = response.UserResponse{
		Code:    http.StatusCreated,
		Message: "User created with success",
		Data:    "",
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, response)
}

func LoginUser(ctx *gin.Context) {

	// create a user request to parse values
	request := request.UserRequestLogin{}

	// validate the parse
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Erro to login access"})
		return
	}

	// search user on database
	var user = response.DataUser{}

	initializers.DB.Raw("SELECT * FROM users WHERE email = ?", request.Email).Scan(&user)

	// verifiy if user exists
	if user.ID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Email or password invalid"})
		return
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Error to compare passwords"})
		return
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_ACCESS")))

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "failed to create token"})
		return
	}
	//set token in cookie
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{})
}

func Validate(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "I'm is logged in"})
}
