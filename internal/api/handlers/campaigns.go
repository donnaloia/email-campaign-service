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

// Get handles GET requests to retrieve a single campaign
func (h *CampaignHandler) Get(c echo.Context) error {
	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Get the campaign ID from the URL
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	// Get the resource
	campaign, err := h.campaignService.GetByID(organizationID, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, campaign)
}

// ListCampaigns handles GET requests to retrieve campaigns
func (h *CampaignHandler) List(c echo.Context) error {
	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Parse pagination parameters from query string
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	// Create pagination params with defaults
	params := models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	// Pass the params to GetAll
	result, err := h.campaignService.GetAll(organizationID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// CreateCampaign handles POST requests to create new campaigns
func (h *CampaignHandler) Create(c echo.Context) error {
	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Bind the request body to the CreateCampaign struct
	var req models.CreateCampaign
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Add organizationID to the request
	req.OrganizationID = organizationID

	// Create the resource
	campaign, err := h.campaignService.Create(organizationID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, campaign)
}

func (h *CampaignHandler) Update(c echo.Context) error {

	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Get the campaign ID from the URL

	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	// Bind the request body to the UpdateCampaign struct
	var req models.UpdateCampaign
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update the resource
	campaign, err := h.campaignService.Update(organizationID, id, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, campaign)
}
