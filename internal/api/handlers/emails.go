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
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	email, err := h.emailService.GetByID(id)
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

	// Pass the params to GetAll
	result, err := h.emailService.GetAll(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// CreateEmail handles POST requests to create new email addresses
func (h *EmailHandler) Create(c echo.Context) error {
	var req models.CreateEmailAddressRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	email, err := h.emailService.Create(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, email)
}
