package handlers

import (
	"net/http"

	"github.com/donnaloia/sendpulse/pkg/response"

	"github.com/labstack/echo/v4"
)

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, response.New(
		"success",
		"Service is healthy",
		nil,
	))
}
