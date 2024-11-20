package models

import "time"

// EmailAddress is a single email address
type EmailAddress struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

// Create a single email address
type CreateEmailAddressRequest struct {
	Address string `json:"address"`
}

// EmailGroup is a group of email addresses
type EmailGroup struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Create a single email group
type CreateEmailGroupRequest struct {
	Name     string   `json:"name"`
	EmailIDs []string `json:"email_ids"` // Array of email address IDs instead of full objects
}

// EmailGroupMember represents the junction between EmailGroup and EmailAddress
type EmailGroupMember struct {
	ID             string    `json:"id"`
	EmailGroupID   string    `json:"email_group_id"`
	EmailAddressID string    `json:"email_address_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Create an email group member
type CreateEmailGroupMember struct {
	EmailGroupID   string `json:"email_group_id"`
	EmailAddressID string `json:"email_address_id"`
}

// Campaign is a high-level object representing a campaign
type Campaign struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Create a single campaign
type CreateCampaignRequest struct {
	Name string `json:"name"`
}

// EmailGroupCampaign is an intermediary model that links an email group to a campaign
type EmailGroupCampaign struct {
	ID           string    `json:"id"`
	EmailGroupID string    `json:"email_group_id"`
	CampaignID   string    `json:"campaign_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// Create an email group campaign
type CreateEmailGroupCampaign struct {
	EmailGroupID string `json:"email_group_id"`
	CampaignID   string `json:"campaign_id"`
}

// Organization is a high-level object representing an Organization
type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Create an organization
type CreateOrganization struct {
	Name string `json:"name"`
}
