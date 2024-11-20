package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/donnaloia/sendpulse/internal/models"
	"github.com/donnaloia/sendpulse/internal/services"

	"github.com/labstack/echo/v4"
)

// Campaigns handler group - capitalized to make it public
var Campaigns *CampaignHandler

// Initialize the campaigns handler
func InitCampaigns(db *sql.DB) {
	Campaigns = &CampaignHandler{
		campaignService: services.NewCampaignService(db),
	}
}

type CampaignHandler struct {
	campaignService *services.CampaignService
}

// GetCampaign handles GET requests to retrieve a single campaign
func (h *CampaignHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	campaign, err := h.campaignService.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, campaign)
}

// ListCampaigns handles GET requests to retrieve campaigns
func (h *CampaignHandler) List(c echo.Context) error {
	// Parse pagination parameters from query string
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	// Create pagination params with defaults
	params := models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	// Pass the params to GetAll
	result, err := h.campaignService.GetAll(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// CreateCampaign handles POST requests to create new campaigns
func (h *CampaignHandler) Create(c echo.Context) error {
	var req models.CreateCampaignRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	campaign, err := h.campaignService.Create(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, campaign)
}
