package services

import (
	"database/sql"
	"fmt"

	"github.com/donnaloia/sendpulse/internal/models"
)

type EmailGroupService struct {
	db *sql.DB
}

func NewEmailGroupService(db *sql.DB) *EmailGroupService {
	return &EmailGroupService{db: db}
}

func (s *EmailGroupService) GetAll(organizationID string, params models.PaginationParams) (*models.PaginatedResponse[models.EmailGroup], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	// Get total count
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM email_groups").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting email groups: %w", err)
	}

	// Calculate offset
	offset := (params.Page - 1) * params.PageSize
	rows, err := s.db.Query(
		`SELECT id, name, organization_id, created_at 
		FROM email_groups 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`,
		params.PageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching email groups: %w", err)
	}
	defer rows.Close()

	var groups []models.EmailGroup
	for rows.Next() {
		var group models.EmailGroup
		if err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.OrganizationID,
			&group.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning email group: %w", err)
		}
		groups = append(groups, group)
	}

	return models.NewPaginatedResponse(groups, total, params.Page, params.PageSize), nil
}

func (s *EmailGroupService) GetByID(organizationID string, id string) (*models.EmailGroup, error) {
	var group models.EmailGroup
	err := s.db.QueryRow(
		"SELECT id, name, organization_id, created_at FROM email_groups WHERE id = $1",
		id,
	).Scan(
		&group.ID,
		&group.Name,
		&group.OrganizationID,
		&group.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("email group not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching email group: %w", err)
	}
	return &group, nil
}

func (s *EmailGroupService) Create(organizationID string, req *models.CreateEmailGroup) (*models.EmailGroup, error) {
	var group models.EmailGroup
	err := s.db.QueryRow(
		`INSERT INTO email_groups (name, organization_id) 
		VALUES ($1, $2) 
		RETURNING id, name, organization_id, created_at`,
		req.Name,
		organizationID,
	).Scan(
		&group.ID,
		&group.Name,
		&group.OrganizationID,
		&group.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating email group: %w", err)
	}
	return &group, nil
}
