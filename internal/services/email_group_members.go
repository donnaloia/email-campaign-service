package services

import (
	"database/sql"
	"fmt"

	"github.com/donnaloia/sendpulse/internal/models"
)

type EmailGroupMemberService struct {
	db *sql.DB
}

func NewEmailGroupMemberService(db *sql.DB) *EmailGroupMemberService {
	return &EmailGroupMemberService{db: db}
}

func (s *EmailGroupMemberService) GetAll(params models.PaginationParams) (*models.PaginatedResponse[models.EmailGroupMember], error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM email_group_members").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error counting group members: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	rows, err := s.db.Query(
		`SELECT id, email_group_id, email_address_id, created_at
		FROM email_group_members
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		params.PageSize, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching group members: %w", err)
	}
	defer rows.Close()

	var members []models.EmailGroupMember
	for rows.Next() {
		var member models.EmailGroupMember
		if err := rows.Scan(
			&member.ID,
			&member.EmailGroupID,
			&member.EmailAddressID,
			&member.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning group member: %w", err)
		}
		members = append(members, member)
	}

	return models.NewPaginatedResponse(members, total, params.Page, params.PageSize), nil
}

func (s *EmailGroupMemberService) GetByID(id string) (*models.EmailGroupMember, error) {
	var member models.EmailGroupMember
	err := s.db.QueryRow(
		`SELECT id, email_group_id, email_address_id, created_at
		FROM email_group_members
		WHERE id = $1`,
		id,
	).Scan(
		&member.ID,
		&member.EmailGroupID,
		&member.EmailAddressID,
		&member.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("group member not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching group member: %w", err)
	}
	return &member, nil
}

func (s *EmailGroupMemberService) Create(req *models.CreateEmailGroupMember) (*models.EmailGroupMember, error) {
	var member models.EmailGroupMember
	err := s.db.QueryRow(
		`INSERT INTO email_group_members (email_group_id, email_address_id)
		VALUES ($1, $2)
		RETURNING id, email_group_id, email_address_id, created_at`,
		req.EmailGroupID, req.EmailAddressID,
	).Scan(
		&member.ID,
		&member.EmailGroupID,
		&member.EmailAddressID,
		&member.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating group member: %w", err)
	}
	return &member, nil
}

func (s *EmailGroupMemberService) Delete(groupID, memberID string) error {
	result, err := s.db.Exec(
		"DELETE FROM email_group_members WHERE email_group_id = $1 AND id = $2",
		groupID, memberID,
	)
	if err != nil {
		return fmt.Errorf("error deleting group member: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deletion result: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("group member not found")
	}

	return nil
}
