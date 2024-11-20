-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email addresses table
CREATE TABLE IF NOT EXISTS email_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    address VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email groups table
CREATE TABLE IF NOT EXISTS email_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email group members (junction table for email_groups and email_addresses)
CREATE TABLE IF NOT EXISTS email_group_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email_group_id UUID NOT NULL REFERENCES email_groups(id) ON DELETE CASCADE,
    email_address_id UUID NOT NULL REFERENCES email_addresses(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(email_group_id, email_address_id)
);

-- Campaigns table
CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email group campaigns (junction table for email_groups and campaigns)
CREATE TABLE IF NOT EXISTS email_group_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email_group_id UUID NOT NULL REFERENCES email_groups(id) ON DELETE CASCADE,
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(email_group_id, campaign_id)
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_email_addresses_address ON email_addresses(address);
CREATE INDEX IF NOT EXISTS idx_email_groups_name ON email_groups(name);
CREATE INDEX IF NOT EXISTS idx_campaigns_name ON campaigns(name);
CREATE INDEX IF NOT EXISTS idx_email_group_members_email_group_id ON email_group_members(email_group_id);
CREATE INDEX IF NOT EXISTS idx_email_group_members_email_address_id ON email_group_members(email_address_id);
CREATE INDEX IF NOT EXISTS idx_email_group_campaigns_email_group_id ON email_group_campaigns(email_group_id);
CREATE INDEX IF NOT EXISTS idx_email_group_campaigns_campaign_id ON email_group_campaigns(campaign_id); 


-- Insert sample data
INSERT INTO organizations (name) VALUES 
    ('sendPulse');

INSERT INTO email_addresses (address) VALUES 
    ('test1@example.com'),
    ('tagl8qosmmbj1x@gmail.com'),
    ('casper.castaneda@yahoo.com'),
    ('9dktv07nvuf@hotmail.com'),
    ('levi.berry@hotmail.com'),
    ('foo04bcnj5hxmkr1v@comcast.net'),
    ('kyaan.rosario@outlook.com'),
    ('s3fffxqsen7ad@msn.com'),
    ('neil.owens@outlook.com'),
    ('vyhmjevatcpk7@yahoo.com'),
    ('abraham.george@yahoo.com'),
    ('ttm4n9p19b5@aol.com'),
    ('test2@example.com'),
    ('paul_stein@bluehost.com');

INSERT INTO email_groups (name) VALUES 
    ('Artists'),
    ('Musicians'),
    ('Songwriters'),
    ('Investors'),
    ('Over 65'),
    ('Under 25'),
    ('Songwriters'),
    ('Impulsive Buyers'),
    ('Everyone'),
    ('Non-subscribers'),
    ('Subscribers'),
    ('Cancelled Subscribers'),
    ('Longterm Subscribers');

INSERT INTO campaigns (name) VALUES 
    ('Summer 2024'),
    ('Winter 2024'),
    ('Spring Into Saving'),
    ('Back to School'),
    ('Holiday 2024'),
    ('Winter 2025'),
    ('Valentine''s Day 2025'),
    ('Easter 2025'),
    ('Summer 2025'),
    ('Special Offers'),
    ('Anniversary'),
    ('Birthday'),
    ('Seasonal'),
    ('New Release'),
    ('Catalog'),
    ('Re-engagement');

-- Create email group memberships by referencing the inserted records
INSERT INTO email_group_members (email_group_id, email_address_id) 
SELECT 
    eg.id as email_group_id,
    ea.id as email_address_id
FROM email_groups eg
CROSS JOIN email_addresses ea
WHERE eg.name = 'Re-engagement' 
    AND ea.address IN ('test1@example.com', 'test2@example.com');

-- Add test2@example.com to Newsletter group
INSERT INTO email_group_members (email_group_id, email_address_id)
SELECT 
    eg.id,
    ea.id
FROM email_groups eg
CROSS JOIN email_addresses ea
WHERE eg.name = 'Longterm Subscribers' 
    AND ea.address IN ('casper.castaneda@yahoo.com', 'test2@example.com', 'foo04bcnj5hxmkr1v@comcast.net');