package middleware

import (
	"jobsy-api/services"
	"jobsy-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OwnershipMiddleware checks if the logged-in user owns the job they're trying to update or delete
func OwnershipMiddleware(jobService *services.JobService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jobID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.Response{
				Status:  "error",
				Message: "Invalid job ID",
			})
			return
		}
		userId, exists := ctx.Get("userId") // Extract the user ID from the JWT claims

		if !exists {
			ctx.JSON(http.StatusUnauthorized, utils.Response{
				Status:  "error",
				Message: "User not authenticated",
			})
			ctx.Abort()
			return
		}

		// Fetch the job by ID
		job, err := jobService.GetJobByID(jobID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Job not found",
			})
			ctx.Abort()
			return
		}

		// Check if the job's company ID matches the logged-in user's company ID
		if job.Company.Hex() != userId {
			ctx.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "You are not authorized to perform this action on this job",
			})
			ctx.Abort()
			return
		}

		// Ownership verified, allow the request to proceed
		ctx.Next()
	}
}
