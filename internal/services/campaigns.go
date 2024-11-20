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

func (s *CampaignService) GetAll(params models.PaginationParams) (*models.PaginatedResponse[models.Campaign], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM campaigns").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting campaigns: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	rows, err := s.db.Query(
		`SELECT id, name, created_at 
		FROM campaigns 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`,
		params.PageSize, offset,
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
			&campaign.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning campaign: %w", err)
		}
		campaigns = append(campaigns, campaign)
	}

	return models.NewPaginatedResponse(campaigns, total, params.Page, params.PageSize), nil
}

func (s *CampaignService) GetByID(id string) (*models.Campaign, error) {
	var campaign models.Campaign
	err := s.db.QueryRow(
		`SELECT id, name, created_at 
		FROM campaigns WHERE id = $1`,
		id,
	).Scan(
		&campaign.ID,
		&campaign.Name,
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

func (s *CampaignService) Create(req *models.CreateCampaignRequest) (*models.Campaign, error) {
	var campaign models.Campaign
	err := s.db.QueryRow(
		`INSERT INTO campaigns (name) 
		VALUES ($1) 
		RETURNING id, name,  created_at`,
		req.Name, "draft",
	).Scan(
		&campaign.ID,
		&campaign.Name,
		&campaign.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating campaign: %w", err)
	}
	return &campaign, nil
}
