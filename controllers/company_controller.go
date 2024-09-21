// func (c *CompanyController) UpdateCompany(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	var company models.Company
// 	if err := ctx.ShouldBindJSON(&company); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	updatedCompany, err := c.companyService.UpdateCompany(id, &company)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, updatedCompany)
// }

package controllers

import (
	"jobsy-api/models"
	"jobsy-api/services"
	"jobsy-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyController struct {
	companyService *services.CompanyService
}

func NewCompanyController(companyService *services.CompanyService) *CompanyController {
	return &CompanyController{companyService: companyService}
}

func (cc *CompanyController) CreateCompany(c *gin.Context) {
	userID, ok := utils.ExtractUserID(c)
	if !ok {
		return
	}
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	company.UserId = userID
	result, err := cc.companyService.CreateCompany(&company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Status:  utils.Success,
		Message: "On boarded successfully",
		Data:    result.InsertedID,
	})
}

func (cc *CompanyController) GetCompanyByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}
	company, err := cc.companyService.GetCompanyByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Company not found"})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Status: utils.Success,
		Data:   company,
	})
}
