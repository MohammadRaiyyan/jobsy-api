package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string             `json:"email" binding:"required,email"`
	Password  string             `json:"password" binding:"required,min=6"`
	Role      string             `bson:"role" json:"role"` // "applicant", "company"
	Token     string             `bson:"token,omitempty" json:"token"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Custom errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrExpiredToken       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid token")
)
