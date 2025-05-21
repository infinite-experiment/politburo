package api

import (
	"encoding/json"
	"infinite-experiment/politburo/internal/models/entities"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

// HealthCheckHandler handles GET /healthCheck
//
// @Summary Health check
// @Description Verifies the server is running.
// @Tags Misc
// @Success 200 {string} string "ok"
// @Router /healthCheck [get]
func HealthCheckHandler(db *sqlx.DB, upSince time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		services := make(map[string]entities.ServiceStatus)

		// Check postgres
		pgstatus := "ok"
		pgDetails := "Postgres Connected"
		if err := db.Ping(); err != nil {
			pgstatus = "down"
			pgDetails = err.Error()
		}
		services["postgres"] = entities.ServiceStatus{
			Status:  pgstatus,
			Details: pgDetails,
		}

		overallStatus := "ok"
		for _, svc := range services {
			if svc.Status != "ok" {
				overallStatus = "down"
				break
			}
		}

		now := time.Now()
		uptime := now.Sub(upSince).Round(time.Second).String()

		resp := entities.HealthCheckResponse{
			Services: services,
			Status:   overallStatus,
			Uptime:   uptime,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
