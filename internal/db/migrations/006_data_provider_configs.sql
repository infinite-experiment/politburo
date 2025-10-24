-- =====================================================
-- Migration 006: Data Provider Configs
-- Purpose: Schema-driven configuration for external data providers (Airtable, etc.)
-- Date: 2025-10-20
-- =====================================================

-- Create ENUM for validation status
DO $$ BEGIN
    CREATE TYPE validation_status AS ENUM ('pending', 'validating', 'valid', 'invalid');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create main configuration table with JSONB for flexible schemas
CREATE TABLE IF NOT EXISTS va_data_provider_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    va_id UUID NOT NULL REFERENCES virtual_airlines(id) ON DELETE CASCADE,

    -- Provider identification
    provider_type VARCHAR(50) NOT NULL,

    -- Configuration data stored as JSONB for flexibility
    config_data JSONB NOT NULL,
    config_version INTEGER NOT NULL DEFAULT 1,

    -- Status flags
    is_active BOOLEAN NOT NULL DEFAULT false,
    validation_status validation_status NOT NULL DEFAULT 'pending',

    -- Enabled features (array of feature names)
    features_enabled TEXT[] DEFAULT '{}',

    -- Validation tracking
    last_validated_at TIMESTAMP,
    validation_errors JSONB,

    -- Audit fields
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Ensure one config per provider type per VA
    UNIQUE(va_id, provider_type)
);

-- Indexes for performance
CREATE INDEX idx_va_provider_configs_va_id ON va_data_provider_configs(va_id);
CREATE INDEX idx_va_provider_configs_active ON va_data_provider_configs(va_id, is_active);
CREATE INDEX idx_va_provider_configs_provider ON va_data_provider_configs(provider_type);

-- GIN index for querying inside JSONB config_data
CREATE INDEX idx_va_provider_configs_data ON va_data_provider_configs USING GIN (config_data);

-- GIN index for validation_errors JSONB
CREATE INDEX idx_va_provider_configs_errors ON va_data_provider_configs USING GIN (validation_errors);

-- Index for features array
CREATE INDEX idx_va_provider_configs_features ON va_data_provider_configs USING GIN (features_enabled);

-- =====================================================
-- Validation History Table
-- Purpose: Track validation results over time for auditing
-- =====================================================

CREATE TABLE IF NOT EXISTS va_provider_validation_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_id UUID NOT NULL REFERENCES va_data_provider_configs(id) ON DELETE CASCADE,

    -- Validation result
    validation_status validation_status NOT NULL,
    validation_errors JSONB,

    -- Validation metadata
    phases_completed TEXT[],
    phases_failed TEXT[],
    duration_ms INTEGER,

    -- Audit
    validated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    triggered_by VARCHAR(50)
);

-- Indexes for validation history table
CREATE INDEX idx_validation_history_config ON va_provider_validation_history(config_id);
CREATE INDEX idx_validation_history_date ON va_provider_validation_history(validated_at);
