package services

import (
	"context"
	"errors"
	"jobsy-api/models"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobService struct {
	jobCollection *mongo.Collection
}

func NewJobService(db *mongo.Client, collection string) *JobService {
	return &JobService{
		jobCollection: db.Database("jobsy-api").Collection(collection),
	}
}

func (s *JobService) CreateJob(job *models.Job) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	job.ID = primitive.NewObjectID()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
	job.Applicants = 0
	return s.jobCollection.InsertOne(ctx, job)
}

func (s *JobService) GetJobByID(id primitive.ObjectID) (*models.Job, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var job models.Job
	err := s.jobCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	return &job, err
}

func (s *JobService) GetAllJobs() ([]*models.Job, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var jobs []*models.Job
	cursor, err := s.jobCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *JobService) UpdateJob(id string, job *models.Job) (*models.Job, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	objectID, _ := primitive.ObjectIDFromHex(id)
	job.UpdatedAt = time.Now()

	updateData := bson.M{}
	jobValue := reflect.ValueOf(job).Elem()
	jobType := reflect.TypeOf(*job)

	for i := 0; i < jobValue.NumField(); i++ {
		fieldValue := jobValue.Field(i)
		fieldType := jobType.Field(i)

		if !isZero(fieldValue) {
			fieldName := fieldType.Tag.Get("bson")
			if fieldName == "status" {
				status := fieldValue.Interface().(models.JobStatus)
				if !isValidJobStatus(status) {
					return nil, errors.New("invalid job status")
				}
			}
			updateData[fieldType.Tag.Get("bson")] = fieldValue.Interface()
		}
	}

	updateData["updated_at"] = job.UpdatedAt

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updateData}

	_, err := s.jobCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	updatedJob := &models.Job{}
	err = s.jobCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(updatedJob)
	if err != nil {
		return nil, err
	}

	return updatedJob, nil

}

func (s *JobService) GetRecentJobs(limit int64) ([]*models.Job, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	options := options.Find().SetSort(bson.M{"postedAt": -1}).SetLimit(limit)
	var jobs []*models.Job
	cursor, err := s.jobCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *JobService) DeleteJob(id string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID}
	_, err := s.jobCollection.DeleteOne(ctx, filter)
	return err
}

func (s *JobService) GetRecommendedJobs(jobID primitive.ObjectID) ([]*models.Job, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"tags": bson.M{"$in": job.Tags},
		"_id":  bson.M{"$ne": job.ID},
	}
	options := options.Find().SetLimit(5)
	var recommendedJobs []*models.Job
	cursor, err := s.jobCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &recommendedJobs)
	if err != nil {
		return nil, err
	}

	return recommendedJobs, nil
}

func (s *JobService) GetJobsByCompany(companyID primitive.ObjectID) ([]*models.Job, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.M{"company": companyID}
	var jobs []*models.Job
	cursor, err := s.jobCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	}
}

func isValidJobStatus(status models.JobStatus) bool {
	switch status {
	case models.Open, models.Closed, models.OnHold, models.Draft:
		return true
	default:
		return false
	}
}
