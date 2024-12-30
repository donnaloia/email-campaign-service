package services

import (
	"database/sql"
	"fmt"

	"github.com/donnaloia/sendpulse/internal/models"
)

type CampaignService struct {
	db *sql.DB
}

func NewCampaignService(db *sql.DB) *CampaignService {
	return &CampaignService{db: db}
}

func (s *CampaignService) GetAll(organizationID string, params models.PaginationParams) (*models.PaginatedResponse[models.Campaign], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	var total int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM campaigns WHERE organization_id = $1",
		organizationID,
	).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting campaigns: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	rows, err := s.db.Query(
		`SELECT id, name, organization_id, created_at 
		FROM campaigns 
		WHERE organization_id = $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`,
		organizationID, params.PageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var campaign models.Campaign
		if err := rows.Scan(
			&campaign.ID,
			&campaign.Name,
			&campaign.OrganizationID,
			&campaign.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning campaign: %w", err)
		}
		campaigns = append(campaigns, campaign)
	}

	return models.NewPaginatedResponse(campaigns, total, params.Page, params.PageSize), nil
}

func (s *CampaignService) GetByID(organizationID string, id string) (*models.Campaign, error) {
	var campaign models.Campaign
	err := s.db.QueryRow(
		`SELECT id, name, organization_id, created_at 
		FROM campaigns 
		WHERE id = $1 AND organization_id = $2`,
		id, organizationID,
	).Scan(
		&campaign.ID,
		&campaign.Name,
		&campaign.OrganizationID,
		&campaign.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("campaign not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching campaign: %w", err)
	}

	campaign.Templates = []models.Template{}

	rows, err := s.db.Query(`
		SELECT t.id, t.name, t.organization_id, t.html, t.created_at
		FROM templates t
		JOIN campaign_templates ct ON ct.template_id = t.id
		WHERE ct.campaign_id = $1`,
		id,
	)

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
		campaign.Templates = append(campaign.Templates, template)
	}

	return &campaign, nil
}

func (s *CampaignService) Create(organizationID string, req *models.CreateCampaign) (*models.Campaign, error) {
	var campaign models.Campaign
	err := s.db.QueryRow(
		`INSERT INTO campaigns (name, organization_id) 
		VALUES ($1, $2) 
		RETURNING id, name, organization_id, created_at`,
		req.Name, organizationID,
	).Scan(
		&campaign.ID,
		&campaign.Name,
		&campaign.OrganizationID,
		&campaign.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating campaign: %w", err)
	}
	return &campaign, nil
}

func (s *CampaignService) Update(organizationID string, id string, req *models.UpdateCampaign) (*models.Campaign, error) {
	// Start a transaction since we're updating multiple tables
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if we don't commit

	// First check if campaign exists and belongs to organization
	_, err = s.GetByID(organizationID, id)
	if err != nil {
		return nil, err
	}

	// Update campaign name if provided
	if req.Name != "" {
		_, err = tx.Exec(
			`UPDATE campaigns 
			 SET name = $1
			 WHERE id = $2 AND organization_id = $3`,
			req.Name, id, organizationID,
		)
		if err != nil {
			return nil, fmt.Errorf("error updating campaign name: %w", err)
		}
	}

	// Update templates if provided
	if req.Templates != nil {
		// First, remove all existing template associations
		_, err = tx.Exec(
			`DELETE FROM campaign_templates 
			 WHERE campaign_id = $1`,
			id,
		)
		if err != nil {
			return nil, fmt.Errorf("error removing existing templates: %w", err)
		}

		// Then add new template associations
		for _, templateID := range req.Templates {
			_, err = tx.Exec(
				`INSERT INTO campaign_templates (campaign_id, template_id)
				 VALUES ($1, $2)`,
				id, templateID,
			)
			if err != nil {
				return nil, fmt.Errorf("error adding template %s: %w", templateID, err)
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// Return the updated campaign
	return s.GetByID(organizationID, id)
}
