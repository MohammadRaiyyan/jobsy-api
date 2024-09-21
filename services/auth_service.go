package services

import (
	"context"
	"jobsy-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthService struct {
	userCollection *mongo.Collection
}

func NewAuthService(db *mongo.Client, collection string) *AuthService {
	return &AuthService{
		userCollection: db.Database("jobsy-api").Collection(collection),
	}
}

func (s *AuthService) FindUserByEmail(email string) (*models.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	err := s.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (s *AuthService) CreateUser(user *models.User) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return s.userCollection.InsertOne(ctx, user)
}

func (s *AuthService) UpdateUserToken(userID primitive.ObjectID, token string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"token": token}}

	_, err := s.userCollection.UpdateOne(context.TODO(), filter, update, options.Update())
	return err
}
