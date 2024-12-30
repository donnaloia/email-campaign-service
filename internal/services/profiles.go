package services

import (
	"database/sql"
	"fmt"

	"github.com/donnaloia/sendpulse/internal/models"
)

type ProfileService struct {
	db *sql.DB
}

func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{db: db}
}

func (s *ProfileService) GetAll(organizationID string, params models.PaginationParams) (*models.PaginatedResponse[models.Profile], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	var total int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM profiles WHERE organization_id = $1",
		organizationID,
	).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting profiles: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	rows, err := s.db.Query(
		`SELECT id, username, email, organization_id, picture_url, created_at 
		FROM profiles 
		WHERE organization_id = $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`,
		organizationID, params.PageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching profiles: %w", err)
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var profile models.Profile
		if err := rows.Scan(
			&profile.ID,
			&profile.Username,
			&profile.Email,
			&profile.OrganizationID,
			&profile.PictureURL,
			&profile.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning profile: %w", err)
		}
		profiles = append(profiles, profile)
	}

	return models.NewPaginatedResponse(profiles, total, params.Page, params.PageSize), nil
}

func (s *ProfileService) GetByID(organizationID string, id string) (*models.Profile, error) {
	var profile models.Profile
	err := s.db.QueryRow(
		`SELECT id, username, email, organization_id, picture_url, created_at 
		FROM profiles 
		WHERE id = $1 AND organization_id = $2`,
		id, organizationID,
	).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.OrganizationID,
		&profile.PictureURL,
		&profile.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching profile: %w", err)
	}
	return &profile, nil
}

func (s *ProfileService) Create(organizationID string, req *models.CreateProfile) (*models.Profile, error) {
	var profile models.Profile
	err := s.db.QueryRow(
		`INSERT INTO profiles (id, username, email, organization_id, picture_url) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, username, email, organization_id, picture_url, created_at`,
		req.ID, req.Username, req.Email, organizationID, req.PictureURL,
	).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.OrganizationID,
		&profile.PictureURL,
		&profile.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating profile: %w", err)
	}
	return &profile, nil
}
