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
var Profiles *ProfileHandler

// Initialize the profiles handler
func InitProfiles(db *sql.DB) {
	Profiles = &ProfileHandler{
		profileService: services.NewProfileService(db),
	}
}

type ProfileHandler struct {
	profileService *services.ProfileService
}

// Get handles GET requests to retrieve a single campaign
func (h *ProfileHandler) Get(c echo.Context) error {
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
	profile, err := h.profileService.GetByID(organizationID, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, profile)
}

// ListCampaigns handles GET requests to retrieve campaigns
func (h *ProfileHandler) List(c echo.Context) error {
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
	result, err := h.profileService.GetAll(organizationID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// CreateProfile handles POST requests to create new profiles
func (h *ProfileHandler) Create(c echo.Context) error {
	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Bind the request body to the CreateProfile struct
	var req models.CreateProfile
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Add organizationID to the request
	req.OrganizationID = organizationID

	// Create the resource
	profile, err := h.profileService.Create(organizationID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, profile)
}
