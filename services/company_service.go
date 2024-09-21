package services

import (
	"context"
	"jobsy-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyService struct {
	companyCollection *mongo.Collection
}

func NewCompanyService(db *mongo.Client, collection string) *CompanyService {
	return &CompanyService{
		companyCollection: db.Database("jobsy-api").Collection(collection),
	}
}

func (s *CompanyService) CreateCompany(company *models.Company) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	company.ID = primitive.NewObjectID()
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()
	return s.companyCollection.InsertOne(ctx, company)
}

func (s *CompanyService) GetCompanyByID(id primitive.ObjectID) (*models.Company, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var company models.Company
	err := s.companyCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&company)
	return &company, err
}

// func (s *CompanyService) UpdateCompany(id string, company *models.Company) (*models.Company, error) {
// 	objectID, _ := primitive.ObjectIDFromHex(id)

// 	filter := bson.M{"_id": objectID}
// 	update := bson.M{"$set": company}
// 	_, err := s.collection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return company, nil
// }

// func (s *CompanyService) GetCompanyByID(id string) (*models.Company, error) {
// 	objectID, _ := primitive.ObjectIDFromHex(id)
// 	filter := bson.M{"_id": objectID}
// 	var company models.Company
// 	err := s.collection.FindOne(context.Background(), filter).Decode(&company)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &company, nil
// }

// func (s *CompanyService) GetCompanyByEmail(email string) (*models.Company, error) {
// 	filter := bson.M{"email": email}
// 	var company models.Company

// 	err := s.collection.FindOne(context.Background(), filter).Decode(&company)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, nil // Return nil if no company is found
// 		}
// 		// Log the error for debugging
// 		return nil, err // Return other errors as is
// 	}
// 	return &company, nil
// }
