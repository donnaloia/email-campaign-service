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
