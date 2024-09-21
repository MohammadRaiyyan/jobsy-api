package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationStatus string

const (
	Pending     ApplicationStatus = "Pending"
	UnderReview ApplicationStatus = "Under Review"
	Hold        ApplicationStatus = "On Hold"
	Rejected    ApplicationStatus = "Rejected"
	Scheduled   ApplicationStatus = "Interview Scheduled"
)

type Applicant struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	JobID               primitive.ObjectID `bson:"jobId" json:"jobId"`
	ApplicantID         primitive.ObjectID `bson:"applicantId" json:"applicantId"`
	Firstname           string             `bson:"firstname" json:"firstname"`
	Lastname            string             `bson:"lastname" json:"lastname"`
	Email               string             `bson:"email" json:"email"`
	Phone               string             `bson:"phone" json:"phone"`
	LinkedinProfile     string             `bson:"linkedinProfile" json:"linkedinProfile"`
	WebsiteURL          string             `bson:"websiteUrl" json:"websiteUrl"`
	Experience          string             `bson:"experience" json:"experience"`
	CurrentJobTitle     string             `bson:"currentJobTitle" json:"currentJobTitle"`
	CurrentEmployerName string             `bson:"currentEmployerName" json:"currentEmployerName"`
	DesiredSalary       string             `bson:"desiredSalary" json:"desiredSalary"`
	AvailabilityDate    string             `bson:"availabilityDate" json:"availabilityDate"`
	Summary             string             `bson:"summary" json:"summary"`
	Resume              string             `bson:"resume" json:"resume"`
	Status              ApplicationStatus  `bson:"status" json:"status"`
	CreatedAt           time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time          `bson:"updatedAt" json:"updatedAt"`
}
