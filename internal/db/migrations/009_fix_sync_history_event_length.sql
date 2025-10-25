-- Increase event column size in va_sync_history to accommodate longer event names
-- PILOT_AT_SYNC, ROUTES_AT_SYNC, PIREPS_AT_SYNC are all 14 characters
ALTER TABLE va_sync_history
ALTER COLUMN event TYPE VARCHAR(50);
