package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	context "infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/services"

	"github.com/go-chi/chi/v5"
)

// UserFlightsHandler godoc
// @Summary      Get flight history for a user
// @Description  Returns a paginated list of flights for the given user.
// @Tags         Flights
// @Produce      json
// @Param        user_id       path     string  true   "User ID"
// @Param        page          query    int     false  "Page number"                   default(1)
// @Param        X-Discord-Id  header   string  true   "Discord ID"                    default(987654321)
// @Param        X-Server-Id   header   string  true   "Discord Server ID"             default(123456789)
// @Param        X-API-Key     header   string  true   "API KEY"                       default(ABCDEF0123456789)
// @Success      200           {object} dtos.APIResponse
// @Failure      400,500       {object} dtos.APIResponse
// @Router       /api/v1/user/{user_id}/flights [get]
func UserFlightsHandler(fltSvc *services.FlightsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		// Extract path parameter
		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			resp := dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      "Invalid IFC ID Received",
				ResponseTime: common.GetResponseTime(initTime),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Parse query parameter 'page'
		page := 1
		if qs := r.URL.Query().Get("page"); qs != "" {
			if p, err := strconv.Atoi(qs); err == nil && p > 0 {
				page = p
			} else {
				resp := dtos.APIResponse{
					Status:       string(constants.APIStatusError),
					Message:      "Invalid page parameter",
					ResponseTime: common.GetResponseTime(initTime),
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}

		// Call service
		dto, err := fltSvc.GetUserFlights(userID, page)
		if err != nil {
			resp := dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      err.Error(),
				ResponseTime: common.GetResponseTime(initTime),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := dtos.APIResponse{
			Status:       string(constants.APIStatusOk),
			Message:      "Fetched Results",
			ResponseTime: common.GetResponseTime(initTime),
			Data:         dto,
		}
		// Send JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func VaFlightsHandler(fltSvc *services.FlightsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()
		claims := context.GetUserClaims(r.Context())

		f, err := fltSvc.GetVALiveFlights(r.Context(), claims.ServerID())

		msg := "Live flights fetched"
		if err != nil {
			msg = err.Error()
		}
		resp := dtos.APIResponse{
			Status:       string(constants.APIStatusError),
			Message:      msg,
			ResponseTime: common.GetResponseTime(initTime),
			Data:         f,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func LiveServers(fltSvc *services.FlightsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		servers, err := fltSvc.GetLiveServers()

		msg := ""
		if err != nil {
			msg = err.Error()
		}
		resp := dtos.APIResponse{
			Status:       string(constants.APIStatusError),
			Message:      msg,
			ResponseTime: common.GetResponseTime(initTime),
			Data:         servers,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
