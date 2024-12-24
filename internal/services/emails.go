package services

import (
	"database/sql"
	"fmt"

	"github.com/donnaloia/sendpulse/internal/models"
)

type EmailService struct {
	db *sql.DB
}

func NewEmailService(db *sql.DB) *EmailService {
	return &EmailService{db: db}
}

func (s *EmailService) GetAll(organizationID string, params models.PaginationParams) (*models.PaginatedResponse[models.EmailAddress], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM email_addresses WHERE organization_id = $1", organizationID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting emails: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	rows, err := s.db.Query(
		`SELECT id, address, organization_id, created_at 
		FROM email_addresses 
		WHERE organization_id = $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`,
		organizationID, params.PageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching emails: %w", err)
	}
	defer rows.Close()

	var emails []models.EmailAddress
	for rows.Next() {
		var email models.EmailAddress
		if err := rows.Scan(
			&email.ID,
			&email.Address,
			&email.OrganizationID,
			&email.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning email: %w", err)
		}
		emails = append(emails, email)
	}

	return models.NewPaginatedResponse(emails, total, params.Page, params.PageSize), nil
}

func (s *EmailService) GetByID(organizationID string, id string) (*models.EmailAddress, error) {
	var email models.EmailAddress
	err := s.db.QueryRow(
		`SELECT id, address, organization_id, created_at 
		FROM email_addresses 
		WHERE id = $1 AND organization_id = $2`,
		id, organizationID,
	).Scan(
		&email.ID,
		&email.Address,
		&email.OrganizationID,
		&email.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("email not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching email: %w", err)
	}
	return &email, nil
}

func (s *EmailService) Create(organizationID string, req *models.CreateEmailAddressRequest) (*models.EmailAddress, error) {
	var email models.EmailAddress
	err := s.db.QueryRow(
		`INSERT INTO email_addresses (address, organization_id) 
		VALUES ($1, $2) 
		RETURNING id, address, organization_id, created_at`,
		req.Address, organizationID,
	).Scan(
		&email.ID,
		&email.Address,
		&email.OrganizationID,
		&email.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating email: %w", err)
	}
	return &email, nil
}
