package routes

import (
	"github.com/donnaloia/sendpulse/internal/api/handlers"

	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo) {
	// Health routes
	e.GET("/health", handlers.HealthCheck)

	// API group
	api := e.Group("/api")

	// Organization Routes
	organizations := api.Group("/organizations")
	organizations.GET("", handlers.Organizations.List)
	organizations.GET("/:id", handlers.Organizations.Get)
	organizations.POST("", handlers.Organizations.Create)

	// Email Routes
	emails := api.Group("/emails")
	emails.GET("", handlers.Emails.List)
	emails.GET("/:id", handlers.Emails.Get)
	emails.POST("", handlers.Emails.Create)

	// Email Group Routes
	emailGroups := api.Group("/email-groups")
	emailGroups.GET("", handlers.EmailGroups.List)
	emailGroups.GET("/:id", handlers.EmailGroups.Get)
	emailGroups.POST("", handlers.EmailGroups.Create)

	// Email Groups Routes
	emailGroupMembers := api.Group("/email-group-members")
	emailGroupMembers.GET("", handlers.EmailGroupMembers.List)
	emailGroupMembers.GET("/:id", handlers.EmailGroupMembers.Get)
	emailGroupMembers.POST("", handlers.EmailGroupMembers.Create)

	// Campaign Routes
	campaigns := api.Group("/campaigns")
	campaigns.GET("", handlers.Campaigns.List)
	emails.GET("/:id", handlers.Campaigns.Get)
	emails.POST("", handlers.Campaigns.Create)
}
