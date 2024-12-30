package routes

import (
	"github.com/donnaloia/sendpulse/internal/api/handlers"

	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo) {
	// Health routes
	e.GET("/health", handlers.HealthCheck)

	// API group
	api := e.Group("/api/v1")

	// Organization as the base route
	organizations := api.Group("/organizations")
	organizations.GET("", handlers.Organizations.List)
	organizations.GET("/:id", handlers.Organizations.Get)
	organizations.POST("", handlers.Organizations.Create)

	// Resources under organization
	org := organizations.Group("/:organization_id")

	// Profile Routes
	profiles := org.Group("/profiles")
	profiles.GET("", handlers.Profiles.List)
	profiles.GET("/:id", handlers.Profiles.Get)
	profiles.POST("", handlers.Profiles.Create)

	// Email Routes
	emails := org.Group("/email-addresses")
	emails.GET("", handlers.Emails.List)
	emails.GET("/:id", handlers.Emails.Get)
	emails.POST("", handlers.Emails.Create)

	// Email Group Routes
	emailGroups := org.Group("/email-groups")
	emailGroups.GET("", handlers.EmailGroups.List)
	emailGroups.GET("/:id", handlers.EmailGroups.Get)
	emailGroups.POST("", handlers.EmailGroups.Create)

	// Email Group Members Routes
	emailGroupMembers := org.Group("/email-group-members")
	emailGroupMembers.GET("", handlers.EmailGroupMembers.List)
	emailGroupMembers.GET("/:id", handlers.EmailGroupMembers.Get)
	emailGroupMembers.POST("", handlers.EmailGroupMembers.Create)

	// Campaign Routes
	campaigns := org.Group("/campaigns")
	campaigns.GET("", handlers.Campaigns.List)
	campaigns.GET("/:id", handlers.Campaigns.Get)
	campaigns.POST("", handlers.Campaigns.Create)
	campaigns.PATCH("/:id", handlers.Campaigns.Update)

	// Template Routes
	templates := org.Group("/templates")
	templates.GET("", handlers.Templates.List)
	templates.GET("/:id", handlers.Templates.Get)
	templates.POST("", handlers.Templates.Create)
}
