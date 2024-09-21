package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId      primitive.ObjectID `bson:"userId" json:"userId"`
	Name        string             `bson:"name" json:"name"`
	CompanySize string             `bson:"companySize" json:"companySize"`
	Avatar      string             `bson:"avatar" json:"avatar"`
	Address     string             `bson:"address" json:"address"`
	City        string             `bson:"city" json:"city"`
	State       string             `bson:"state" json:"state"`
	Pincode     int                `bson:"pincode" json:"pincode"`
	Country     string             `bson:"country" json:"country"`
	Phone       string             `bson:"phone" json:"phone"`
	Website     string             `bson:"website" json:"website"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
