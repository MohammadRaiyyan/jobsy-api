package middleware

import (
	"errors"
	"jobsy-api/models"
	"jobsy-api/services"
	"jobsy-api/utils"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	bearerSchema = "Bearer "
)

func AuthMiddleware(authService *services.AuthService, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response{Status: utils.Error, Message: err.Error()})
			c.Abort()
			return
		}

		claims, err := validateToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response{Status: utils.Error, Message: "Invalid authorization token"})
			c.Abort()
			return
		}

		user, err := authenticateUser(claims, authService)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response{Status: utils.Error, Message: err.Error()})
			c.Abort()
			return
		}

		if user.Token != token {
			c.JSON(http.StatusUnauthorized, utils.Response{Status: utils.Error, Message: "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("email", user.Email)
		c.Set("role", user.Role)
		c.Set("userId", user.ID)
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}
	if !strings.HasPrefix(authHeader, bearerSchema) {
		return "", errors.New("authorization header must start with 'Bearer'")
	}
	return strings.TrimPrefix(authHeader, bearerSchema), nil
}

func validateToken(tokenString, jwtSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if err := validateTokenExpiration(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func validateTokenExpiration(claims jwt.MapClaims) error {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid expiration claim")
	}
	if time.Now().Unix() > int64(exp) {
		return errors.New("token has expired")
	}
	return nil
}

func authenticateUser(claims jwt.MapClaims, authService *services.AuthService) (*models.User, error) {
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid user")
	}

	user, err := authService.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
