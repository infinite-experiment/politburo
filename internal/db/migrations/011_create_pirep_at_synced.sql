-- Create PIREP sync table to store flight logs from Airtable

CREATE TABLE IF NOT EXISTS pirep_at_synced (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    at_id VARCHAR(20) NOT NULL,
    server_id UUID NOT NULL,

    -- Core PIREP fields
    route TEXT,                          -- Route name/identifier
    flight_mode VARCHAR(50),             -- "Career Mode", "General", "Training", etc.
    flight_time NUMERIC(10, 2),          -- Flight duration in hours (e.g., 2.50)
    pilot_callsign VARCHAR(50),          -- Pilot callsign
    aircraft VARCHAR(100),               -- Aircraft type/model
    livery VARCHAR(100),                 -- Airline/livery

    -- References to synced records (Airtable IDs)
    route_at_id VARCHAR(20),             -- Reference to route_at_synced.at_id
    pilot_at_id VARCHAR(20),             -- Reference to pilot_at_synced.at_id

    -- Airtable metadata
    at_created_time TIMESTAMP,           -- When the PIREP was created in Airtable (for sorting)

    -- Timestamps
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),

    -- Unique constraint
    CONSTRAINT pirep_at_synced_unique UNIQUE (server_id, at_id)
);

-- Indexes for performance
CREATE INDEX idx_pirep_at_synced_server_id ON pirep_at_synced(server_id);
CREATE INDEX idx_pirep_at_synced_pilot_callsign ON pirep_at_synced(server_id, pilot_callsign);
CREATE INDEX idx_pirep_at_synced_pilot_at_id ON pirep_at_synced(server_id, pilot_at_id);
CREATE INDEX idx_pirep_at_synced_route_at_id ON pirep_at_synced(server_id, route_at_id);
CREATE INDEX idx_pirep_at_synced_flight_mode ON pirep_at_synced(server_id, flight_mode);
CREATE INDEX idx_pirep_at_synced_at_created_time ON pirep_at_synced(server_id, at_created_time DESC);
