-- Create VA Config Table
CREATE TABLE va_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    va_id UUID REFERENCES virtual_airlines(id),
    config_key VARCHAR(50) NOT NULL,
    config_value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (va_id, config_key)
);

-- Indexes for performance
CREATE INDEX idx_va_configs_va_id ON va_configs(va_id);
CREATE INDEX idx_va_configs_va_key ON va_configs(va_id, config_key);

ALTER TABLE public.virtual_airlines
    DROP COLUMN IF EXISTS callsign_prefix,
    DROP COLUMN IF EXISTS callsign_suffix;