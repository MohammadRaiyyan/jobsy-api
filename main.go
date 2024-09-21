package main

import (
	"context"
	"jobsy-api/controllers"
	"jobsy-api/routes"
	"jobsy-api/services"
	"log"
	"os"

	"jobsy-api/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file available")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatalln("mongodb uri string not found : ")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if uri == "" {
		log.Fatalln("JWT token string not found : ")
	}
	client, err := config.ConnectDB(uri)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
			return
		}
	}()

	jobService := services.NewJobService(client, "jobs")
	companyService := services.NewCompanyService(client, "companies")
	applicantService := services.NewApplicantService(client, "applicants", "jobs")
	authService := services.NewAuthService(client, "users")

	jobController := controllers.NewJobController(jobService)
	companyController := controllers.NewCompanyController(companyService)
	applicantController := controllers.NewApplicantController(applicantService)
	authController := controllers.NewAuthController(authService, jwtSecret)

	routes.SetupRoutes(router, jobController, companyController, applicantController, authController, jobService, authService, jwtSecret)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
