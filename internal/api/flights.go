package api

import (
	"encoding/json"
	"fmt"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/services"
	"net/http"
	"strconv"
	"time"
)

// InitUserRegistrationHandler handles GET /api/v1/user/{user_id}/flights
//
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
func UserFlightsHandler(fltScv *services.FlightsService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v\n========\n%v", r, *r)
		initTime := time.Now()
		user_id := r.PathValue("user_id")

		if user_id == "" {
			resp := dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      "Invalid IFC ID Received",
				ResponseTime: services.GetResponseTime(initTime),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		page := 1

		qs := r.URL.Query().Get("page")
		if qs != "" {
			if num, err := strconv.Atoi(qs); err == nil && num > 0 {
				page = num
			}
		}

		flights, err := fltScv.GetUserFlights(user_id, page)

		if err != nil {
			resp := dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      err.Error(),
				ResponseTime: services.GetResponseTime(initTime),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
		}

		resp := dtos.APIResponse{
			Status:       string(constants.APIStatusOk),
			Message:      "Success",
			Data:         flights,
			ResponseTime: services.GetResponseTime(initTime),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
