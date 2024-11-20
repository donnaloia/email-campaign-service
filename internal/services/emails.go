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

func (s *EmailService) GetAll(params models.PaginationParams) (*models.PaginatedResponse[models.EmailAddress], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM email_addresses").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting emails: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	rows, err := s.db.Query(
		`SELECT id, address, created_at 
		FROM email_addresses 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`,
		params.PageSize, offset,
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
			&email.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning email: %w", err)
		}
		emails = append(emails, email)
	}

	return models.NewPaginatedResponse(emails, total, params.Page, params.PageSize), nil
}

func (s *EmailService) GetByID(id string) (*models.EmailAddress, error) {
	var email models.EmailAddress
	err := s.db.QueryRow(
		"SELECT id, address, created_at FROM email_addresses WHERE id = $1",
		id,
	).Scan(
		&email.ID,
		&email.Address,
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

func (s *EmailService) Create(req *models.CreateEmailAddressRequest) (*models.EmailAddress, error) {
	var email models.EmailAddress
	err := s.db.QueryRow(
		`INSERT INTO email_addresses (address) 
		VALUES ($1) 
		RETURNING id, address, created_at`,
		req.Address,
	).Scan(
		&email.ID,
		&email.Address,
		&email.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating email: %w", err)
	}
	return &email, nil
}
