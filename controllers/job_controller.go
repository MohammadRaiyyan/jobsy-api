package controllers

import (
	"jobsy-api/models"
	"jobsy-api/services"
	"jobsy-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobController struct {
	jobService *services.JobService
}

func NewJobController(jobService *services.JobService) *JobController {
	return &JobController{jobService: jobService}
}

func (jc *JobController) CreateJob(c *gin.Context) {
	userID, ok := utils.ExtractUserID(c)
	if !ok {
		return
	}

	// Bind the request body to the Job model
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign the ObjectID to the job.Company field
	job.Company = userID

	// Create the job using the job service
	result, err := jc.jobService.CreateJob(&job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, result)
}

func (jc *JobController) GetJobByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}
	job, err := jc.jobService.GetJobByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (c *JobController) GetAllJobs(ctx *gin.Context) {
	jobs, err := c.jobService.GetAllJobs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Status:  utils.Error,
			Message: "Failed to retrieve job listings",
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status:  utils.Success,
		Message: "Job listings retrieved successfully",
		Data:    jobs,
	})
}

func (c *JobController) UpdateJob(ctx *gin.Context) {
	id := ctx.Param("id")
	var job models.Job
	if err := ctx.ShouldBindJSON(&job); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Status:  utils.Error,
			Message: "Invalid job update payload",
		})
		return
	}

	updatedJob, err := c.jobService.UpdateJob(id, &job)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Status:  utils.Error,
			Message: "Failed to update the job",
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status:  utils.Success,
		Message: "Job updated successfully",
		Data:    updatedJob,
	})
}

func (c *JobController) GetRecentPostings(ctx *gin.Context) {
	jobs, err := c.jobService.GetRecentJobs(10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Status:  utils.Error,
			Message: "Failed to retrieve recent job postings",
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status:  utils.Error,
		Message: "Recent job postings retrieved successfully",
		Data:    jobs,
	})
}

func (c *JobController) DeleteJob(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.jobService.DeleteJob(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Status:  utils.Error,
			Message: "Failed to delete the job",
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status:  utils.Success,
		Message: "Job deleted successfully",
	})
}

func (c *JobController) GetRecommendedJobs(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Status:  utils.Error,
			Message: "Invalid job ID",
		})
		return
	}

	jobs, err := c.jobService.GetRecommendedJobs(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Status:  utils.Error,
			Message: "Failed to retrieve recommended jobs",
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Recommended jobs retrieved successfully",
		Data:    jobs,
	})
}

func (c *JobController) GetJobsByCompany(ctx *gin.Context) {
	userId, ok := utils.ExtractUserID(ctx)

	if !ok {
		return
	}

	jobs, err := c.jobService.GetJobsByCompany(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Status:  "error",
			Message: "Failed to retrieve jobs",
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Jobs for the company retrieved successfully",
		Data:    jobs,
	})
}
