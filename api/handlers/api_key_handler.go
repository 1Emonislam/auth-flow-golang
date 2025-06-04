package handlers

import (
	"net/http"
	"own-paynet/api/response"
	"own-paynet/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type APIKeyHandler struct {
	apiKeyService *services.APIKeyService
}

func NewAPIKeyHandler(apiKeyService *services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

type GenerateAPIKeyRequest struct {
	Description string `json:"description" binding:"required"`
}

func (h *APIKeyHandler) GenerateAPIKey(c *gin.Context) {
	var req GenerateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	apiKey, err := h.apiKeyService.GenerateAPIKey(userID, req.Description)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "API key generated successfully", apiKey)
}

func (h *APIKeyHandler) GetUserAPIKeys(c *gin.Context) {
	userID := c.GetUint("user_id")
	apiKeys, err := h.apiKeyService.GetUserAPIKeys(userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "API keys retrieved successfully", apiKeys)
}

func (h *APIKeyHandler) SetDefaultKey(c *gin.Context) {
	userID := c.GetUint("user_id")
	apiKeyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid API key ID")
		return
	}

	err = h.apiKeyService.SetDefaultKey(userID, uint(apiKeyID))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Default API key updated successfully", nil)
}

func (h *APIKeyHandler) DeleteAPIKey(c *gin.Context) {
	userID := c.GetUint("user_id")
	apiKeyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid API key ID")
		return
	}

	err = h.apiKeyService.DeleteAPIKey(userID, uint(apiKeyID))
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "API key deleted successfully", nil)
}
