package models

import "time"

const (
	CampaignStatusDraft     = "draft"
	CampaignStatusScheduled = "scheduled"
	CampaignStatusLaunched  = "launched"
)

// EmailAddress is a single email address
type EmailAddress struct {
	ID             string    `json:"id"`
	Address        string    `json:"address"`
	OrganizationID string    `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Create a single email address
type CreateEmailAddressRequest struct {
	Address        string `json:"address"`
	OrganizationID string `json:"organization_id"`
}

// EmailGroup is a group of email addresses
type EmailGroup struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationID string    `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Create a single email group
type CreateEmailGroup struct {
	Name           string   `json:"name"`
	OrganizationID string   `json:"organization_id"`
	EmailIDs       []string `json:"email_ids"` // Array of email address IDs instead of full objects
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
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Status         string       `json:"status"`
	OrganizationID string       `json:"organization_id"`
	CreatedAt      time.Time    `json:"created_at"`
	Templates      []Template   `json:"templates"`
	EmailGroups    []EmailGroup `json:"email_groups"`
}

// Create a single campaign
type CreateCampaign struct {
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id"`
}

// Update a single campaign
type UpdateCampaign struct {
	Name        string   `json:"name,omitempty"`
	Status      string   `json:"status,omitempty"`
	Templates   []string `json:"templates,omitempty"`
	EmailGroups []string `json:"email_groups,omitempty"` // Array of email group IDs
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

// Template is a high-level object representing a template
type Template struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	HTML           string    `json:"html"`
	CreatedAt      time.Time `json:"created_at"`
}

// Create a template
type CreateTemplate struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	HTML           string `json:"html"`
}

// CampaignTemplate is an intermediary model that links a campaign to a template
type CampaignTemplate struct {
	ID         string    `json:"id"`
	CampaignID string    `json:"campaign_id"`
	TemplateID string    `json:"template_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Create a campaign template
type CreateCampaignTemplate struct {
	CampaignID string `json:"campaign_id"`
	TemplateID string `json:"template_id"`
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

// Profile is a high-level object representing a Profile
type Profile struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	OrganizationID string    `json:"organization_id"`
	PictureURL     string    `json:"picture_url"`
	CreatedAt      time.Time `json:"created_at"`
}

// Create a profile
type CreateProfile struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	OrganizationID string `json:"organization_id"`
	PictureURL     string `json:"picture_url"`
}
