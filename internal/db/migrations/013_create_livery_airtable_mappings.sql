-- Migration: 013_create_livery_airtable_mappings.sql
-- Purpose: Create table for livery-to-Airtable field mappings (aircraft/airline standardization)

-- Create table for livery airtable mappings
CREATE TABLE IF NOT EXISTS livery_airtable_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    va_id UUID NOT NULL,
    livery_id VARCHAR(255) NOT NULL,
    field_type VARCHAR(50) NOT NULL, -- 'aircraft' or 'airline'
    source_value VARCHAR(255) NOT NULL, -- Raw from IF API
    target_value VARCHAR(255) NOT NULL, -- Standardized for Airtable
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(va_id, livery_id, field_type),
    FOREIGN KEY (va_id) REFERENCES virtual_airlines(id) ON DELETE CASCADE
);

-- Create indexes for query performance
CREATE INDEX IF NOT EXISTS idx_livery_mappings_va_livery ON livery_airtable_mappings(va_id, livery_id);
CREATE INDEX IF NOT EXISTS idx_livery_mappings_lookup ON livery_airtable_mappings(va_id, field_type, source_value);
CREATE INDEX IF NOT EXISTS idx_livery_mappings_va_id ON livery_airtable_mappings(va_id);
