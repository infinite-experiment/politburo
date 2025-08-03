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

ALTER TABLE virtual_airlines
    DROP COLUMN IF EXISTS callsign_prefix,
    DROP COLUMN IF EXISTS callsign_suffix;

ALTER TABLE virtual_airlines
ADD COLUMN IF NOT EXISTS is_airtable_enabled boolean DEFAULT false;

ALTER TABLE va_user_roles
ADD COLUMN IF NOT EXISTS airtable_pilot_id character varying(20),
ADD COLUMN IF NOT EXISTS callsign character varying(20)
    

CREATE TABLE IF NOT EXISTS public.va_sync_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    va_id UUID NOT NULL,
    event VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
    last_sync_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS public.pilot_at_synced (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    at_id VARCHAR(20) NOT NULL,
    callsign VARCHAR(20),
    registered BOOLEAN DEFAULT false
);

ALTER TABLE pilot_at_synced
ADD COLUMN server_id UUID;


	ALTER TABLE pilot_at_synced
ADD CONSTRAINT pilot_at_synced_server_at_id_key UNIQUE (server_id, at_id);


CREATE TABLE IF NOT EXISTS route_at_synced (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    at_id VARCHAR(20) NOT NULL,
    server_id UUID NOT NULL,
    origin VARCHAR(10),
    destination VARCHAR(10),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),

    CONSTRAINT route_at_synced_unique UNIQUE (server_id, at_id)
);

ALTER TABLE route_at_synced
ADD COLUMN IF NOT EXISTS route VARCHAR(20);

