package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobStatus string

const (
	Open   JobStatus = "Open"
	Closed JobStatus = "Closed"
	OnHold JobStatus = "On Hold"
	Draft  JobStatus = "Draft"
)

type WorkType string

const (
	Remote     WorkType = "Remote"
	HybridWork WorkType = "Hybrid Work"
	OnSite     WorkType = "On Site"
)

type JobType string

const (
	FullTime   JobType = "Full Time"
	Contract   JobType = "Contract"
	Temporary  JobType = "Temporary"
	Internship JobType = "Internship"
	Fresher    JobType = "Fresher"
)

type Job struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Company     primitive.ObjectID `bson:"company" json:"company"`
	Location    string             `bson:"location" json:"location"`
	JobType     JobType            `bson:"jobType" json:"jobType"`
	WorkType    WorkType           `bson:"workType" json:"workType"`
	Salary      Salary             `bson:"salary" json:"salary"`
	Summary     string             `bson:"summary" json:"summary"`
	Description string             `bson:"description" json:"description"`
	Applicants  int                `bson:"applicants" json:"applicants"`
	Status      JobStatus          `bson:"status" json:"status"`
	Tags        []string           `bson:"tags" json:"tags"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Salary struct {
	Negotiable   bool   `bson:"negotiable" json:"negotiable"`
	Min          string `bson:"min,omitempty" json:"min,omitempty"`
	Max          string `bson:"max,omitempty" json:"max,omitempty"`
	CurrencyType string `bson:"currencyType,omitempty" json:"currencyType,omitempty"`
}
