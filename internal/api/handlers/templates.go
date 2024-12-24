package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/donnaloia/sendpulse/internal/models"
	"github.com/donnaloia/sendpulse/internal/services"
	"github.com/labstack/echo/v4"
)

var Templates *TemplateHandler

func InitTemplates(db *sql.DB) {
	Templates = &TemplateHandler{
		templateService: services.NewTemplateService(db),
	}
}

type TemplateHandler struct {
	templateService *services.TemplateService
}

func (h *TemplateHandler) Get(c echo.Context) error {
	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Get the template ID from the URL
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	// Get the resource
	template, err := h.templateService.GetByID(organizationID, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, template)
}

func (h *TemplateHandler) List(c echo.Context) error {
	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Get the pagination parameters from the query string
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	// Create the pagination parameters
	params := models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	// Get all the resources
	result, err := h.templateService.GetAll(organizationID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TemplateHandler) Create(c echo.Context) error {
	// Get the organization ID from the URL
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing organization ID")
	}

	// Bind the request body to the CreateTemplateRequest struct
	var req models.CreateTemplate
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create the resource
	template, err := h.templateService.Create(organizationID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, template)
}
