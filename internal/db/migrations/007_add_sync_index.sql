-- =====================================================================
-- Migration 007: Add Index for Incremental Sync Performance
-- =====================================================================
-- Purpose: Optimize the MAX(updated_at) query used for incremental sync
-- Date: 2025-10-20
-- =====================================================================

-- Index to speed up finding the most recent sync timestamp per VA
-- This supports the query: SELECT MAX(updated_at) FROM va_user_roles WHERE va_id = ? AND airtable_pilot_id IS NOT NULL
CREATE INDEX IF NOT EXISTS idx_va_user_roles_sync_timestamp
    ON va_user_roles(va_id, updated_at DESC)
    WHERE airtable_pilot_id IS NOT NULL;

COMMENT ON INDEX idx_va_user_roles_sync_timestamp IS 'Optimizes incremental sync by efficiently finding last sync timestamp per VA';
