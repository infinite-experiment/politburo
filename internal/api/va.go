package api

import (
	"encoding/json"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	ctxutil "infinite-experiment/politburo/internal/context"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/services"
	"log"
	"net/http"
	"strconv"
	"time"
)

func InitRegisterServer(regService *services.RegistrationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		// ── 1. Bind JSON ─────────────────────────────────────────────
		var req dtos.InitServerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// ── 2. Pull claims (optional) & LOG everything ──────────────
		claims := ctxutil.GetUserClaims(r.Context())

		log.Printf("[InitServer] guild=%s user=%s  ▶  va_code=%q prefix=%q suffix=%q",
			claims.ServerID(), claims.UserID(), req.VACode, req.Prefix, req.Suffix)

		status, steps, err := regService.InitServerRegistration(r.Context(), req.VACode, req.Prefix, req.Suffix, req.VAName)

		if err != nil {
			resp := dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      "Failed to process",
				ResponseTime: common.GetResponseTime(initTime),
				Data: dtos.InitServerResponse{
					VACode:  req.VACode,
					Status:  false,
					Message: "Failed to insert server",
					Steps:   steps,
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := dtos.APIResponse{
			Status:       strconv.FormatBool(status),
			ResponseTime: common.GetResponseTime(initTime),
			Data: dtos.InitServerResponse{
				VACode:  req.VACode,
				Status:  true,
				Message: "Request Processed Successfully",
				Steps:   steps,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
