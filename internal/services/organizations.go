package services

import (
	"database/sql"
	"fmt"

	"github.com/donnaloia/sendpulse/internal/models"
)

type OrganizationService struct {
	db *sql.DB
}

func NewOrganizationService(db *sql.DB) *OrganizationService {
	return &OrganizationService{db: db}
}

func (s *OrganizationService) GetAll(params models.PaginationParams) (*models.PaginatedResponse[models.Organization], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM organizations").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting organizations: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	rows, err := s.db.Query(
		`SELECT id, name, created_at 
		FROM organizations 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`,
		params.PageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching organizations: %w", err)
	}
	defer rows.Close()

	var orgs []models.Organization
	for rows.Next() {
		var org models.Organization
		if err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning organization: %w", err)
		}
		orgs = append(orgs, org)
	}

	return models.NewPaginatedResponse(orgs, total, params.Page, params.PageSize), nil
}

func (s *OrganizationService) GetByID(id string) (*models.Organization, error) {
	var org models.Organization
	err := s.db.QueryRow(
		"SELECT id, name, created_at FROM organizations WHERE id = $1",
		id,
	).Scan(
		&org.ID,
		&org.Name,
		&org.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching organization: %w", err)
	}
	return &org, nil
}

func (s *OrganizationService) Create(req *models.CreateOrganization) (*models.Organization, error) {
	var org models.Organization
	err := s.db.QueryRow(
		`INSERT INTO organizations (name) 
		VALUES ($1) 
		RETURNING id, name, created_at`,
		req.Name,
	).Scan(
		&org.ID,
		&org.Name,
		&org.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating organization: %w", err)
	}
	return &org, nil
}
