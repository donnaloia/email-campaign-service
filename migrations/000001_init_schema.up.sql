-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Profiles table
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    timezone VARCHAR(100),
    bio VARCHAR(255),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    picture_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email addresses table
CREATE TABLE IF NOT EXISTS email_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    address VARCHAR(255) NOT NULL UNIQUE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email groups table
CREATE TABLE IF NOT EXISTS email_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
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
    status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'launched')),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_campaign_status CHECK (status IN ('draft', 'scheduled', 'launched'))
);

-- Email group campaigns (junction table for email_groups and campaigns)
CREATE TABLE IF NOT EXISTS email_group_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email_group_id UUID NOT NULL REFERENCES email_groups(id) ON DELETE CASCADE,
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(email_group_id, campaign_id)
);

-- Templates table
CREATE TABLE IF NOT EXISTS templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    html TEXT NOT NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Campaign templates (junction table for campaigns and templates)
CREATE TABLE IF NOT EXISTS campaign_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(campaign_id, template_id)
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_campaigns_organization_id ON campaigns(organization_id);
CREATE INDEX IF NOT EXISTS idx_email_addresses_address ON email_addresses(address);
CREATE INDEX IF NOT EXISTS idx_email_groups_name ON email_groups(name);
CREATE INDEX IF NOT EXISTS idx_campaigns_name ON campaigns(name);
CREATE INDEX IF NOT EXISTS idx_email_group_members_email_group_id ON email_group_members(email_group_id);
CREATE INDEX IF NOT EXISTS idx_email_group_members_email_address_id ON email_group_members(email_address_id);
CREATE INDEX IF NOT EXISTS idx_email_group_campaigns_email_group_id ON email_group_campaigns(email_group_id);
CREATE INDEX IF NOT EXISTS idx_email_group_campaigns_campaign_id ON email_group_campaigns(campaign_id);
CREATE INDEX IF NOT EXISTS idx_templates_name ON templates(name);
CREATE INDEX IF NOT EXISTS idx_templates_organization_id ON templates(organization_id);
CREATE INDEX IF NOT EXISTS idx_profiles_organization_id ON profiles(organization_id);
CREATE INDEX IF NOT EXISTS idx_profiles_email ON profiles(email);


-- Insert sample data
INSERT INTO organizations (name) VALUES 
    ('sendPulse'),
    ('Test Org');;

INSERT INTO email_addresses (address, organization_id) 
SELECT address, (SELECT id FROM organizations WHERE name = 'sendPulse')
FROM (VALUES 
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
    ('paul_stein@bluehost.com')
) AS t(address);

INSERT INTO email_groups (name, organization_id) 
SELECT name, (SELECT id FROM organizations WHERE name = 'sendPulse')
FROM (VALUES 
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
    ('Longterm Subscribers')
) AS t(name);

INSERT INTO campaigns (name, status, organization_id) 
SELECT 
    name,
    status,
    (SELECT id FROM organizations WHERE name = 'sendPulse')
FROM (VALUES
    ('Fall 2025', 'draft'),
    ('Summer 2024', 'launched'),
    ('Winter 2024', 'scheduled'),
    ('Spring Into Saving', 'draft'),
    ('Back to School', 'launched'),
    ('Holiday 2024', 'scheduled'),
    ('Winter 2025', 'draft'),
    ('Valentine''s Day 2025', 'scheduled'),
    ('Easter 2025', 'launched'),
    ('Summer 2025', 'draft'),
    ('Special Offers', 'launched'),
    ('Anniversary', 'scheduled'),
    ('Birthday', 'draft'),
    ('Seasonal', 'launched'),
    ('New Release', 'scheduled'),
    ('Catalog', 'draft'),
    ('Re-engagement', 'launched')
) AS t(name, status);

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

-- Insert template records
INSERT INTO templates (name, html, organization_id)
SELECT name, html, (SELECT id FROM organizations WHERE name = 'sendPulse')
FROM (VALUES 
    ('Welcome Email', '<html><head>
    <title>Your Newsletter Title</title>
    <style>
        body {
            font-family: sans-serif;
            font-size: 16px;
            line-height: 1.5;
            margin: 0;
            padding: 20px;
            background-color: #f4f4f4;
        }

        h1 {
            color: #333;
            font-size: 24px;
            margin-bottom: 20px;
        }

        p {
            color: #555;
            margin-bottom: 15px;
        }

        a {
            color: #007bff;
            text-decoration: none;
        }

        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #fff;
            padding: 20px;
            border-radius: 5px;
        }

        .button {
            background-color: #007bff;
            color: #fff;
            padding: 10px 20px;
            border: none;
            border-radius: 3px;
            text-decoration: none;
            display: inline-block;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to Our Newsletter!</h1>
        <p>This is our first newsletter. Were excited to share some exciting news and updates with you.</p>
        <p><strong>Heres whats new:</strong></p>
        <ul>
            <li>Weve launched a new product!</li>
            <li>Weve updated our pricing plans.</li>
            <li>Were now hiring for several exciting roles.</li>
        </ul>
        <p>Learn more about these updates by visiting our <a href="https://www.yourwebsite.com">website</a>.</p>
        <p>We hope you enjoy this newsletter. Stay tuned for more updates in the future.</p>
        <p><a href="https://www.yourwebsite.com/blog" class="button">Read Our Blog</a></p>
    </div>
</body></html>'),
    ('Monthly Newsletter', '<div><h2>Monthly Updates</h2><p>Here''s what''s new this month...</p><div class="content-area"></div></div>'),
    ('Product Launch', '<div style="background-color: #f8f8f8;"><h1>New Release!</h1><p>Check out our latest product...</p><img src="product.jpg" /></div>'),
    ('Holiday Special', '<div class="festive"><h1>Season''s Greetings!</h1><p>Celebrate with our special offers...</p><div class="offers"></div></div>'),
    ('Birthday Wishes', '<div style="text-align: center;"><h1>🎉 Happy Birthday! 🎂</h1><p>Here''s a special gift for your special day...</p></div>'),
    ('Abandoned Cart', '<div><h2>Still Interested?</h2><p>Your cart is waiting for you...</p><button class="cta">Complete Purchase</button></div>'),
    ('Event Invitation', '<div class="event"><h1>You''re Invited!</h1><p>Join us for an exclusive event...</p><button>RSVP Now</button></div>'),
    ('Feedback Request', '<div><h2>We Value Your Opinion</h2><p>Please take our quick survey...</p><div class="survey-embed"></div></div>'),
    ('Order Confirmation', '<div class="receipt"><h1>Thank You for Your Order</h1><p>Order #: {{order_number}}</p><div class="order-details"></div></div>'),
    ('Password Reset', '<div style="max-width: 600px;"><h2>Reset Your Password</h2><p>Click the link below to reset your password:</p><a href="{{reset_link}}">Reset Password</a></div>'),
    ('Weekly Digest', '<div class="digest"><h1>This Week''s Highlights</h1><ul>{{#each highlights}}<li>{{this}}</li>{{/each}}</ul></div>')
) AS t(name, html);