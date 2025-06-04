package handlers

import (
	"net/http"
	"strconv"

	"own-paynet/api/response"
	"own-paynet/models"
	"own-paynet/services"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	service services.CompanyService
}

func NewCompanyHandler(s services.CompanyService) *CompanyHandler {
	return &CompanyHandler{s}
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid company ID")
		return
	}

	var input models.Company
	if errBad := c.ShouldBindJSON(&input); errBad != nil {
		response.ErrorResponse(c, http.StatusBadRequest, errBad.Error())
		return
	}

	// Call the service to update
	updatedCompany, err := h.service.UpdateCompany(uint(id), &input)
	if err != nil {
		if err.Error() == "record not found" {
			response.ErrorResponse(c, http.StatusNotFound, "Company not found")
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update company")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Company updated successfully", updatedCompany)
}
