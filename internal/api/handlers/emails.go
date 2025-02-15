package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/donnaloia/sendpulse/internal/models"
	"github.com/donnaloia/sendpulse/internal/services"

	"github.com/labstack/echo/v4"
)

// Items handler group - capitalized to make it public
var Emails *EmailHandler

// Initialize the emails handler
func InitEmails(db *sql.DB) {
	Emails = &EmailHandler{
		emailService: services.NewEmailService(db),
	}
}

type EmailHandler struct {
	emailService *services.EmailService
}

// GetEmail handles GET requests to retrieve a single email address
func (h *EmailHandler) Get(c echo.Context) error {
	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Get the email ID from the URL
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	// Get the resource
	email, err := h.emailService.GetByID(id, organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, email)
}

// ListEmails handles GET requests to retrieve email addresses
func (h *EmailHandler) List(c echo.Context) error {
	// Parse pagination parameters from query string
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	// Create pagination params with defaults
	params := models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Get all the resources
	result, err := h.emailService.GetAll(organizationID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// CreateEmail handles POST requests to create new email addresses
func (h *EmailHandler) Create(c echo.Context) error {
	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Bind the request body to the CreateEmailAddressRequest struct
	var req models.CreateEmailAddressRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create the resource
	email, err := h.emailService.Create(organizationID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, email)
}
