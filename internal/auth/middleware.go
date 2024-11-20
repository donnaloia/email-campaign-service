package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthResponse struct {
	UserID      string   `json:"user_id"`
	Permissions []string `json:"permissions"`
}

// Middleware checks permissions for all routes
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip auth for certain paths
			if isPublicPath(c.Path()) {
				return next(c)
			}

			// Get token from Authorization header
			token := extractToken(c.Request())
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization token")
			}

			// Call auth service to validate token and get permissions
			authResp, err := validateTokenWithAuthService(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			// Get required permission for this path and method
			requiredPerm := getRequiredPermission(c.Path(), c.Request().Method)
			if requiredPerm == "" {
				// No permission required for this path
				return next(c)
			}

			// Check if user has required permission
			if !hasPermission(authResp.Permissions, requiredPerm) {
				return echo.NewHTTPError(http.StatusForbidden,
					fmt.Sprintf("insufficient permissions: %s required", requiredPerm))
			}

			// Store auth response in context for later use
			c.Set("auth", authResp)

			return next(c)
		}
	}
}

// isPublicPath checks if the given path should skip authentication
func isPublicPath(path string) bool {
	publicPaths := []string{
		"/health",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh",
	}

	for _, p := range publicPaths {
		if p == path {
			return true
		}
	}
	return false
}

// extractToken gets the JWT token from the Authorization header
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		return bearerToken[7:]
	}
	return ""
}

// validateTokenWithAuthService calls the auth service to validate the token
func validateTokenWithAuthService(token string) (*AuthResponse, error) {
	// This would typically make an HTTP request to your auth service
	// For now, this is a mock implementation
	if token == "" {
		return nil, fmt.Errorf("invalid token")
	}

	// Mock response
	return &AuthResponse{
		UserID:      "user123",
		Permissions: []string{"email:send", "template:read"},
	}, nil
}

// getRequiredPermission returns the required permission for a given path and method
func getRequiredPermission(path string, method string) string {
	// Define your permission mappings here
	permissionMap := map[string]map[string]string{
		"/api/v1/emails": {
			"POST":   "email:send",
			"GET":    "email:read",
			"DELETE": "email:delete",
		},
		"/api/v1/templates": {
			"POST":   "template:create",
			"GET":    "template:read",
			"PUT":    "template:update",
			"DELETE": "template:delete",
		},
	}

	if methodPerms, ok := permissionMap[path]; ok {
		if perm, ok := methodPerms[method]; ok {
			return perm
		}
	}
	return ""
}

// hasPermission checks if the user has the required permission
func hasPermission(userPerms []string, requiredPerm string) bool {
	for _, p := range userPerms {
		if p == requiredPerm {
			return true
		}
	}
	return false
}
