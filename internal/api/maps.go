package api

import (
	"encoding/json"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/dtos"
	"log"
	"net/http"
	"time"
)

// UserFlightMapHandler godoc
// @Summary      Get cached flight route
// @Description  Returns cached flight information for a given flight ID (from query param `i`)
// @Tags         Flights
// @Accept       json
// @Produce      json
// @Param        i   query     string  true  "Flight ID"
// @Success      200 {object} dtos.APIResponse
// @Failure      404 {object} dtos.APIResponse
// @Router       /public/flight [get]
func UserFlightMapHandler(c *common.CacheService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()
		flightID := r.URL.Query().Get("i")

		log.Printf("API CALLED: %s", flightID)
		if flightID == "" {
			resp := &dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      "Missing required flight ID",
				ResponseTime: common.GetResponseTime(initTime),
			}
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		result := common.GetFlightFromCache(c, flightID)

		if result == nil {
			resp := &dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      "Flight details not found or unavailable. Please try to regenerate the link via /logbook command",
				ResponseTime: common.GetResponseTime(initTime),
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := &dtos.APIResponse{
			Status:       string(constants.APIStatusOk),
			Message:      "Data found",
			ResponseTime: common.GetResponseTime(initTime),
			Data:         result,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

	}
}
