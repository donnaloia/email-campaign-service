package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/donnaloia/sendpulse/internal/models"
	"github.com/donnaloia/sendpulse/internal/services"

	"github.com/labstack/echo/v4"
)

// EmailGroups handler group - capitalized to make it public
var EmailGroups *EmailGroupHandler

// Initialize the email groups handler
func InitEmailGroups(db *sql.DB) {
	EmailGroups = &EmailGroupHandler{
		emailGroupService: services.NewEmailGroupService(db),
	}
}

type EmailGroupHandler struct {
	emailGroupService *services.EmailGroupService
}

// GetEmailGroup handles GET requests to retrieve a single email group
func (h *EmailGroupHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	emailGroup, err := h.emailGroupService.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, emailGroup)
}

// ListEmailGroups handles GET requests to retrieve email groups
func (h *EmailGroupHandler) List(c echo.Context) error {
	// Check if service is initialized
	if h.emailGroupService == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "email group service not initialized")
	}

	// Parse pagination parameters from query string
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	// Create pagination params with defaults
	params := models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	// Get paginated email groups from service
	result, err := h.emailGroupService.GetAll(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error fetching email groups: %v", err))
	}

	return c.JSON(http.StatusOK, result)
}

// CreateEmailGroup handles POST requests to create a new email group
func (h *EmailGroupHandler) Create(c echo.Context) error {
	var emailGroup models.EmailGroup
	if err := c.Bind(&emailGroup); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO: Implement your database insertion logic here
	return c.JSON(http.StatusCreated, emailGroup)
}

// UpdateEmailGroup handles PUT requests to update an existing email group
func (h *EmailGroupHandler) UpdateEmailGroup(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	var emailGroup models.EmailGroup
	if err := c.Bind(&emailGroup); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO: Implement your database update logic here
	return c.JSON(http.StatusOK, emailGroup)
}

// DeleteEmailGroup handles DELETE requests to delete an email group
func (h *EmailGroupHandler) DeleteEmailGroup(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	// TODO: Implement your database deletion logic here
	return c.NoContent(http.StatusNoContent)
}
