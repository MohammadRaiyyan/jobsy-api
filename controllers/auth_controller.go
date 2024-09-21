package controllers

import (
	"jobsy-api/models"
	"jobsy-api/services"
	"jobsy-api/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	authService *services.AuthService
	jwtSecret   string
}

func NewAuthController(authService *services.AuthService, jwtSecret string) *AuthController {
	return &AuthController{authService: authService, jwtSecret: jwtSecret}
}

func (ac *AuthController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Status: utils.Error, Message: err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Status: utils.Error, Message: "Could not hash password"})
		return
	}
	user.Password = string(hashedPassword)

	_, err = ac.authService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Status: utils.Error, Message: "Could not create user"})
		return
	}

	c.JSON(http.StatusOK, utils.Response{Status: utils.Success, Message: "User registered successfully"})
}

func (ac *AuthController) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Status: utils.Error, Message: "Invalid credentials"})
		return
	}

	user, err := ac.authService.FindUserByEmail(credentials.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Response{Status: utils.Error, Message: "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, utils.Response{Status: utils.Error, Message: "Incorrect password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  user.Email,
		"role":   user.Role,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
		"userId": user.ID.Hex(),
	})

	tokenString, err := token.SignedString([]byte(ac.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Status: utils.Error, Message: "Something went wrong while generating unique token"})
		return
	}

	user.Token = tokenString
	if err := ac.authService.UpdateUserToken(user.ID, tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Status: utils.Error, Message: "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (ac *AuthController) Logout(c *gin.Context) {
	userID, ok := utils.ExtractUserID(c)
	if !ok {
		return
	}

	if err := ac.authService.UpdateUserToken(userID, ""); err != nil {
		c.JSON(http.StatusInternalServerError,
			utils.Response{Status: utils.Success, Message: "Something went wrong could not able to logged out"},
		)
		return
	}

	c.JSON(http.StatusOK, utils.Response{Status: utils.Success, Message: "Logged out successfully"})
}
