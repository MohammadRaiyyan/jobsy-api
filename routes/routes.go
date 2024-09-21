package routes

import (
	"jobsy-api/controllers"
	"jobsy-api/middleware"
	"jobsy-api/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	jobController *controllers.JobController,
	companyController *controllers.CompanyController,
	applicantController *controllers.ApplicantController,
	authController *controllers.AuthController,
	jobService *services.JobService,
	authService *services.AuthService,
	jwtSecret string,
) {
	router.POST("/auth/register", authController.Register)
	router.POST("/auth/login", authController.Login)

	public := router.Group("/api")
	public.GET("/jobs", jobController.GetAllJobs)
	public.GET("/jobs/recent", jobController.GetRecentPostings)
	public.GET("/jobs/:id", jobController.GetJobByID)
	public.GET("/jobs/recommended/:id", jobController.GetRecommendedJobs)

	public.POST("/companies", companyController.CreateCompany)

	auth := router.Group("/api")
	auth.Use(middleware.AuthMiddleware(authService, jwtSecret))

	// Jobs Routes
	auth.POST("/jobs", jobController.CreateJob)

	auth.GET("/jobs/company", jobController.GetJobsByCompany)
	auth.PUT("/jobs/:id", middleware.OwnershipMiddleware(jobService), jobController.UpdateJob)
	auth.DELETE("/jobs/:id", middleware.OwnershipMiddleware(jobService), jobController.DeleteJob)

	// Companies Routes
	auth.GET("/companies/:id", companyController.GetCompanyByID)

	// Applicants Routes
	auth.POST("/applicants", applicantController.CreateApplicant)
	auth.GET("/applicants/job/:id", applicantController.GetApplicantsByJobID)
	auth.GET("/applicants/:id", applicantController.GetApplicantByID)
	auth.PUT("/applicants/:id", applicantController.UpdateApplicantStatus)
	auth.POST("/applicants/:id", applicantController.UpdateApplicantStatus)

	// Logout
	auth.POST("/auth/logout", authController.Logout)
}
