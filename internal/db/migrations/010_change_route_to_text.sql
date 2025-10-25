-- Change route field from VARCHAR(20) to TEXT
-- to support route names of any length without limits

ALTER TABLE route_at_synced
ALTER COLUMN route TYPE TEXT;
