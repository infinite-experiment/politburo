package api

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler handles GET /healthCheck
//
// @Summary Health check
// @Description Verifies the server is running.
// @Tags Misc
// @Success 200 {string} string "ok"
// @Router /healthCheck [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Health Ok")
}
