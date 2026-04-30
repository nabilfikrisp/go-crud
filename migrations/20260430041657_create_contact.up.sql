-- Create enum type for relationship
CREATE TYPE contact_relationship AS ENUM (
    'Friend',
    'Family',
    'Colleague',
    'Other'
);

-- Create contacts table
CREATE TABLE contacts (
    id UUID PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    relationship contact_relationship NOT NULL DEFAULT 'Other',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT contacts_email_unique UNIQUE (email)
);

-- Indexes for common filter queries
CREATE INDEX idx_contacts_first_name ON contacts (first_name);

CREATE INDEX idx_contacts_last_name ON contacts (last_name);

CREATE INDEX idx_contacts_email ON contacts (email);

CREATE INDEX idx_contacts_phone_number ON contacts (phone_number);

CREATE INDEX idx_contacts_relationship ON contacts (relationship);

CREATE INDEX idx_contacts_created_at ON contacts (created_at DESC);