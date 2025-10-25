-- Migration: 012_add_flight_modes_config.sql
-- Purpose: Add flight_modes_config JSONB column to virtual_airlines table

-- Add flight_modes_config column to virtual_airlines table
ALTER TABLE virtual_airlines
ADD COLUMN IF NOT EXISTS flight_modes_config JSONB DEFAULT '{}';

-- Add index for JSONB queries
CREATE INDEX IF NOT EXISTS idx_virtual_airlines_flight_modes ON virtual_airlines USING GIN (flight_modes_config);
