package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/donnaloia/sendpulse/internal/models"
	"github.com/donnaloia/sendpulse/internal/services"

	"github.com/labstack/echo/v4"
)

// EmailGroupMembers handler group - capitalized to make it public
var EmailGroupMembers *EmailGroupMemberHandler

// Initialize the email group members handler
func InitEmailGroupMembers(db *sql.DB) {
	EmailGroupMembers = &EmailGroupMemberHandler{
		service: services.NewEmailGroupMemberService(db),
	}
}

type EmailGroupMemberHandler struct {
	service *services.EmailGroupMemberService
}

// GetMember handles GET requests to retrieve a single email group member
func (h *EmailGroupMemberHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing ID")
	}

	member, err := h.service.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, member)
}

// ListMembers handles GET requests to retrieve email group members
func (h *EmailGroupMemberHandler) List(c echo.Context) error {
	// Parse pagination parameters from query string
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

// AddMember handles POST requests to add an email to a group
func (h *EmailGroupMemberHandler) Create(c echo.Context) error {
	groupID := c.Param("group_id")
	if groupID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing group ID")
	}

	var req models.CreateEmailGroupMember
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	req.EmailGroupID = groupID

	member, err := h.service.Create(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, member)
}

// RemoveMember handles DELETE requests to remove an email from a group
func (h *EmailGroupMemberHandler) RemoveMember(c echo.Context) error {
	groupID := c.Param("group_id")
	memberID := c.Param("member_id")
	if groupID == "" || memberID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing group ID or member ID")
	}

	err := h.service.Delete(groupID, memberID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
