package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/donnaloia/sendpulse/internal/models"
	"github.com/donnaloia/sendpulse/internal/services"

	"github.com/labstack/echo/v4"
)

// Organizations handler group - capitalized to make it public
var Organizations *OrganizationHandler

// Initialize the organizations handler
func InitOrganizations(db *sql.DB) {
	Organizations = &OrganizationHandler{
		service: services.NewOrganizationService(db),
	}
}

type OrganizationHandler struct {
	service *services.OrganizationService
}

// Get handles GET requests to retrieve a single organization
func (h *OrganizationHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	org, err := h.service.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, org)
}

// List handles GET requests to retrieve organizations
func (h *OrganizationHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	params := models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.service.GetAll(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// Create handles POST requests to create new organizations
func (h *OrganizationHandler) Create(c echo.Context) error {
	var req models.CreateOrganization
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	org, err := h.service.Create(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, org)
}
