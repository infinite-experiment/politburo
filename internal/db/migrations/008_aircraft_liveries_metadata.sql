-- Migration: 008_aircraft_liveries_metadata.sql
-- Description: Add aircraft_liveries table for persistent storage of IF aircraft/livery metadata
-- Author: System
-- Date: 2025-10-24

-- Create aircraft_liveries table
CREATE TABLE IF NOT EXISTS aircraft_liveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    livery_id VARCHAR(100) UNIQUE NOT NULL,
    aircraft_id VARCHAR(100) NOT NULL,
    aircraft_name TEXT NOT NULL,
    livery_name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_synced_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_aircraft_liveries_livery_id ON aircraft_liveries(livery_id);
CREATE INDEX IF NOT EXISTS idx_aircraft_liveries_aircraft_id ON aircraft_liveries(aircraft_id);
CREATE INDEX IF NOT EXISTS idx_aircraft_liveries_active ON aircraft_liveries(is_active);

-- Composite index for common queries (active liveries by ID)
CREATE INDEX IF NOT EXISTS idx_aircraft_liveries_active_livery ON aircraft_liveries(is_active, livery_id);

-- Index for sync operations
CREATE INDEX IF NOT EXISTS idx_aircraft_liveries_sync ON aircraft_liveries(last_synced_at);

-- Add comment for documentation
COMMENT ON TABLE aircraft_liveries IS 'Stores Infinite Flight aircraft and livery metadata synced from IF API';
COMMENT ON COLUMN aircraft_liveries.livery_id IS 'Infinite Flight livery UUID';
COMMENT ON COLUMN aircraft_liveries.aircraft_id IS 'Infinite Flight aircraft UUID';
COMMENT ON COLUMN aircraft_liveries.is_active IS 'False if livery removed from IF API';
COMMENT ON COLUMN aircraft_liveries.last_synced_at IS 'Last time this record was verified against IF API';
