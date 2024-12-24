package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PermissionsResponse struct {
	UserID      string   `json:"user_id"`
	OrgID       string   `json:"org_id"`
	Permissions []string `json:"permissions"`
}

// Middleware checks permissions for all routes
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// // uncomment this out to skip auth for all routes
			// return next(c)

			// Skip auth for certain paths
			// if isPublicPath(c.Path()) {
			// 	return next(c)
			// }

			// // Get token from Authorization header
			// access_token := extractToken(c.Request())
			// if access_token == "" {
			// 	return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization token")
			// }

			// // Call auth service to validate token and get permissions
			// authResp, err := validateTokenWithAuthService(access_token)
			// if err != nil {
			// 	return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			// }

			// // Get org ID from URL if present (e.g., /api/v1/orgs/:orgId/...)
			// if orgID := c.Param("orgId"); orgID != "" {
			// 	// Check if user belongs to this org
			// 	if orgID != authResp.OrgID {
			// 		return echo.NewHTTPError(http.StatusForbidden, "user does not belong to this organization")
			// 	}
			// }

			// // Get required permission for this path and method
			// requiredPerm := getRequiredPermission(c.Path(), c.Request().Method)
			// if requiredPerm == "" {
			// 	// No permission required for this path
			// 	return next(c)
			// }

			// // Check if user has required permission
			// if !hasPermission(authResp.Permissions, requiredPerm) {
			// 	return echo.NewHTTPError(http.StatusForbidden,
			// 		fmt.Sprintf("insufficient permissions: %s required", requiredPerm))
			// }

			// // Store auth response in context for later use
			// c.Set("auth", authResp)

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
func validateTokenWithAuthService(token string) (*PermissionsResponse, error) {
	// Parse the JWT token to extract the subject (uuid)
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Replace with your actual secret key
		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request to permissions endpoint using Subject from claims
	url := fmt.Sprintf("http://localhost:8080/api/v1/users/%s/permissions", claims.Subject)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling permissions service: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("permissions service returned status: %d", resp.StatusCode)
	}

	// Parse response
	var authResp PermissionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	// Set the UserID from the token subject
	authResp.UserID = claims.Subject

	return &authResp, nil
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
