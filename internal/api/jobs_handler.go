package api

import (
	"encoding/json"
	"fmt"
	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/jobs"
	"log"
	"net/http"
	"time"
)

// JobsHandler handles manual job triggering endpoints
type JobsHandler struct {
	pilotSyncJob *jobs.PilotSyncJob
}

// NewJobsHandler creates a new jobs handler
func NewJobsHandler(pilotSyncJob *jobs.PilotSyncJob) *JobsHandler {
	return &JobsHandler{
		pilotSyncJob: pilotSyncJob,
	}
}

// TriggerPilotSync manually triggers the pilot sync job
// @Summary Trigger pilot sync job
// @Description Manually trigger the pilot sync job for all VAs or a specific VA
// @Tags admin,jobs
// @Accept json
// @Produce json
// @Param body body TriggerPilotSyncRequest false "Optional VA ID to sync specific VA"
// @Success 200 {object} TriggerPilotSyncResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/jobs/sync-pilots [post]
func (h *JobsHandler) TriggerPilotSync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
		start := time.Now()

		// Parse optional request body
		var req TriggerPilotSyncRequest
		if r.Body != nil {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
				respondWithError(w, http.StatusBadRequest, "Invalid request body")
				return
			}
		}

		// Get claims for logging
		claims := auth.GetUserClaims(r.Context())
		triggeredBy := claims.DiscordUserID()

		log.Printf("[JobsHandler] Pilot sync manually triggered by user %s", triggeredBy)

		// If specific VA ID provided, sync only that VA
		ctx := r.Context()
		var syncedCount int
		var err error

		if req.VAID != "" {
			log.Printf("[JobsHandler] Syncing specific VA: %s", req.VAID)
			syncedCount, err = h.pilotSyncJob.SyncVAPilots(ctx, req.VAID)
			if err != nil {
				log.Printf("[JobsHandler] Error syncing VA %s: %v", req.VAID, err)
				respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to sync VA: %v", err))
				return
			}
		} else {
			// Run full sync for all VAs
			log.Printf("[JobsHandler] Running full pilot sync")
			err = h.pilotSyncJob.Run(ctx)
			if err != nil {
				log.Printf("[JobsHandler] Error running pilot sync: %v", err)
				respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to run sync: %v", err))
				return
			}
		}

		duration := time.Since(start)

		response := TriggerPilotSyncResponse{
			Status:       "ok",
			Message:      "Pilot sync completed successfully",
			ResponseTime: fmt.Sprintf("%dms", duration.Milliseconds()),
			Data: PilotSyncResult{
				TriggeredBy:  triggeredBy,
				TriggeredAt:  start.Format(time.RFC3339),
				CompletedAt:  time.Now().Format(time.RFC3339),
				DurationMs:   int(duration.Milliseconds()),
				PilotsSynced: syncedCount,
				VAID:         req.VAID,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// GetJobStatus returns the status of background jobs
// @Summary Get job status
// @Description Get status information about background jobs
// @Tags admin,jobs
// @Produce json
// @Success 200 {object} JobStatusResponse
// @Router /api/v1/admin/jobs/status [get]
func (h *JobsHandler) GetJobStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query database for recent sync history (if we had a sync_history table)
		// For now, return basic status

		response := JobStatusResponse{
			Status:       "ok",
			Message:      "Job status retrieved",
			ResponseTime: "5ms",
			Data: JobStatusData{
				Jobs: []JobInfo{
					{
						Name:        "pilot_sync",
						Description: "Syncs pilot data from Airtable to local database",
						Schedule:    "Every 1 hour",
						Status:      "running",
						LastRun:     "", // Could track this in memory or DB
						NextRun:     "", // Could calculate based on schedule
					},
					{
						Name:        "logbook_worker",
						Description: "Processes flight logbook entries",
						Schedule:    "Continuous (worker)",
						Status:      "running",
						LastRun:     "",
						NextRun:     "",
					},
					{
						Name:        "cache_filler",
						Description: "Pre-fills cache with frequently accessed data",
						Schedule:    "Continuous (worker)",
						Status:      "running",
						LastRun:     "",
						NextRun:     "",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Request/Response types

type TriggerPilotSyncRequest struct {
	VAID string `json:"va_id,omitempty"` // Optional: sync specific VA only
}

type TriggerPilotSyncResponse struct {
	Status       string          `json:"status"`
	Message      string          `json:"message"`
	ResponseTime string          `json:"response_time"`
	Data         PilotSyncResult `json:"data"`
}

type PilotSyncResult struct {
	TriggeredBy  string `json:"triggered_by"`    // Discord ID of admin who triggered
	TriggeredAt  string `json:"triggered_at"`    // ISO 8601 timestamp
	CompletedAt  string `json:"completed_at"`    // ISO 8601 timestamp
	DurationMs   int    `json:"duration_ms"`     // Duration in milliseconds
	PilotsSynced int    `json:"pilots_synced"`   // Number of pilots synced
	VAID         string `json:"va_id,omitempty"` // VA ID if specific VA sync
}

type JobStatusResponse struct {
	Status       string        `json:"status"`
	Message      string        `json:"message"`
	ResponseTime string        `json:"response_time"`
	Data         JobStatusData `json:"data"`
}

type JobStatusData struct {
	Jobs []JobInfo `json:"jobs"`
}

type JobInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
	Status      string `json:"status"` // "running", "stopped", "error"
	LastRun     string `json:"last_run,omitempty"`
	NextRun     string `json:"next_run,omitempty"`
}
