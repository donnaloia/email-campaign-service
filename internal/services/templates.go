package services

import (
	"database/sql"
	"fmt"

	"github.com/donnaloia/sendpulse/internal/models"
)

type TemplateService struct {
	db *sql.DB
}

func NewTemplateService(db *sql.DB) *TemplateService {
	return &TemplateService{db: db}
}

func (s *TemplateService) GetAll(organizationID string, params models.PaginationParams) (*models.PaginatedResponse[models.Template], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	var total int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM templates WHERE organization_id = $1",
		organizationID,
	).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting templates: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	rows, err := s.db.Query(
		`SELECT id, name, organization_id, html, created_at 
		FROM templates 
		WHERE organization_id = $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`,
		organizationID, params.PageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching templates: %w", err)
	}
	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var template models.Template
		if err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.OrganizationID,
			&template.HTML,
			&template.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning template: %w", err)
		}
		templates = append(templates, template)
	}

	return models.NewPaginatedResponse(templates, total, params.Page, params.PageSize), nil
}

func (s *TemplateService) GetByID(organizationID string, id string) (*models.Template, error) {
	var template models.Template
	err := s.db.QueryRow(
		`SELECT id, name, organization_id, html, created_at 
		FROM templates 
		WHERE id = $1 AND organization_id = $2`,
		id, organizationID,
	).Scan(
		&template.ID,
		&template.Name,
		&template.OrganizationID,
		&template.HTML,
		&template.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("template not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching template: %w", err)
	}
	return &template, nil
}

func (s *TemplateService) Create(organizationID string, req *models.CreateTemplate) (*models.Template, error) {
	var template models.Template
	err := s.db.QueryRow(
		`INSERT INTO templates (name, organization_id, html) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, organization_id, html, created_at`,
		req.Name, organizationID, req.HTML,
	).Scan(
		&template.ID,
		&template.Name,
		&template.OrganizationID,
		&template.HTML,
		&template.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating template: %w", err)
	}
	return &template, nil
}
