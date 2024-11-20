package api

import (
	"database/sql"
	"fmt"

	"github.com/donnaloia/sendpulse/internal/api/handlers"
	"github.com/donnaloia/sendpulse/internal/api/middleware"
	"github.com/donnaloia/sendpulse/internal/api/routes"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	db   *sql.DB
}

func NewServer(db *sql.DB) *Server {
	e := echo.New()

	// Verify db connection
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}

	// Initialize handlers with database connection
	handlers.InitEmails(db)
	handlers.InitEmailGroups(db)
	handlers.InitCampaigns(db)
	handlers.InitEmailGroupMembers(db)
	handlers.InitOrganizations(db)

	// Add middleware
	middleware.Setup(e)

	// Setup routes
	routes.Setup(e)

	return &Server{
		echo: e,
		db:   db,
	}
}

func (s *Server) Start(addr string) error {
	return s.echo.Start(addr)
}
