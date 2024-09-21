package services

import (
	"context"
	"errors"
	"jobsy-api/models"
	"net/mail"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicantService struct {
	applicantCollection *mongo.Collection
	jobCollection       *mongo.Collection
}

func NewApplicantService(db *mongo.Client, applicantCollectionName, jobCollectionName string) *ApplicantService {
	return &ApplicantService{
		applicantCollection: db.Database("jobsy-api").Collection(applicantCollectionName),
		jobCollection:       db.Database("jobsy-api").Collection(jobCollectionName),
	}
}

func (s *ApplicantService) CreateApplicant(applicant *models.Applicant, userID primitive.ObjectID) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	applicant.ID = primitive.NewObjectID()
	applicant.ApplicantID = userID
	applicant.CreatedAt = time.Now()
	applicant.UpdatedAt = time.Now()
	applicant.Status = models.Pending

	if err := validateApplicant(applicant); err != nil {
		return err
	}

	_, err := s.applicantCollection.InsertOne(ctx, applicant)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": applicant.JobID}
	update := bson.M{"$inc": bson.M{"applicants": 1}}
	_, err = s.jobCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func validateApplicant(applicant *models.Applicant) error {
	// Validate Job ID
	if applicant.JobID.IsZero() {
		return errors.New("job ID is required")
	}

	// Validate Applicant ID
	if applicant.ApplicantID.IsZero() {
		return errors.New("applicant ID is required")
	}

	// Validate First Name and Last Name
	if applicant.Firstname == "" {
		return errors.New("first name is required")
	}
	if applicant.Lastname == "" {
		return errors.New("last name is required")
	}

	// Validate Email
	if applicant.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(applicant.Email); err != nil {
		return errors.New("invalid email format")
	}

	// Validate Phone Number
	if applicant.Phone != "" {
		phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
		if !phoneRegex.MatchString(applicant.Phone) {
			return errors.New("invalid phone number")
		}
	}

	// Validate LinkedIn Profile (if provided, must be a valid URL)
	if applicant.LinkedinProfile != "" {
		if !isValidURL(applicant.LinkedinProfile) {
			return errors.New("invalid LinkedIn profile URL")
		}
	}

	// Validate Website URL (if provided, must be a valid URL)
	if applicant.WebsiteURL != "" {
		if !isValidURL(applicant.WebsiteURL) {
			return errors.New("invalid website URL")
		}
	}

	// Validate Experience
	if applicant.Experience == "" {
		return errors.New("experience is required")
	}

	// Validate Current Job Title
	if applicant.CurrentJobTitle == "" {
		return errors.New("current job title is required")
	}

	// Validate Current Employer Name
	if applicant.CurrentEmployerName == "" {
		return errors.New("current employer name is required")
	}

	// Validate Desired Salary
	if applicant.DesiredSalary == "" {
		return errors.New("desired salary is required")
	}

	// Validate Availability Date (optional, but must be valid if provided)
	if applicant.AvailabilityDate != "" {
		if !isValidDate(applicant.AvailabilityDate) {
			return errors.New("invalid availability date format")
		}
	}

	// Validate Summary
	if applicant.Summary == "" {
		return errors.New("summary is required")
	}

	// Validate Resume
	if applicant.Resume == "" {
		return errors.New("resume is required")
	}

	// All validations passed
	return nil
}

func isValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^https?://[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`)
	return urlRegex.MatchString(url)
}

func isValidDate(date string) bool {
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return dateRegex.MatchString(date)
}

func (s *ApplicantService) GetApplicantsByJobID(jobID string) ([]*models.Applicant, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, err
	}

	cursor, err := s.applicantCollection.Find(ctx, bson.M{"jobId": objectID})
	if err != nil {
		return nil, err
	}

	var applicants []*models.Applicant
	if err = cursor.All(ctx, &applicants); err != nil {
		return nil, err
	}

	return applicants, nil
}

func (s *ApplicantService) GetApplicantByID(id string) (*models.Applicant, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var applicant models.Applicant
	err = s.applicantCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&applicant)
	if err != nil {
		return nil, err
	}

	return &applicant, nil
}

func (s *ApplicantService) UpdateApplicantStatus(id string, status models.ApplicationStatus) (*models.Applicant, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	_, err = s.applicantCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		return nil, err
	}

	var updatedApplicant models.Applicant
	err = s.applicantCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedApplicant)
	if err != nil {
		return nil, err
	}

	return &updatedApplicant, nil
}
