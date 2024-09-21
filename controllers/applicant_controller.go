package controllers

import (
	"jobsy-api/models"
	"jobsy-api/services"
	"jobsy-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicantController struct {
	applicantService *services.ApplicantService
}

func NewApplicantController(applicantService *services.ApplicantService) *ApplicantController {
	return &ApplicantController{applicantService: applicantService}
}

func (ac *ApplicantController) CreateApplicant(c *gin.Context) {
	userId, ok := utils.ExtractUserID(c)
	if !ok {
		return
	}

	jobID := c.Param("id")

	var applicant models.Applicant
	if err := c.ShouldBindJSON(&applicant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	applicant.JobID, _ = primitive.ObjectIDFromHex(jobID)
	err := ac.applicantService.CreateApplicant(&applicant, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Application successfully submitted",
	})
}

func (c *ApplicantController) GetApplicantsByJobID(ctx *gin.Context) {
	jobID := ctx.Param("id")

	applicants, err := c.applicantService.GetApplicantsByJobID(jobID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Status:  utils.Error,
			Message: err.Error(),
			Data:    applicants,
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status:  utils.Success,
		Message: "",
		Data:    applicants,
	})
}

func (c *ApplicantController) GetApplicantByID(ctx *gin.Context) {
	id := ctx.Param("id")

	applicant, err := c.applicantService.GetApplicantByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.Response{
			Status: utils.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status: utils.Success,
		Data:   applicant,
	})
}

func (c *ApplicantController) UpdateApplicantStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	// Define a structure to bind only the status field
	var statusUpdate struct {
		Status models.ApplicationStatus `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&statusUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{Message: err.Error(), Status: utils.Success})
		return
	}

	// Call the service to update the status
	updatedApplicant, err := c.applicantService.UpdateApplicantStatus(id, statusUpdate.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{Message: err.Error(), Status: utils.Success})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Status: utils.Success,
		Data:   updatedApplicant,
	})
}
