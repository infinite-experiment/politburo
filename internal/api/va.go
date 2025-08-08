package api

import (
	"encoding/json"
	"fmt"
	ctxutil "infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
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

		log.Printf("[InitServer] guild=%s user=%s  ▶  va_code=%q ",
			claims.ServerID(), claims.UserID(), req.VACode)

		status, steps, err := regService.InitServerRegistration(r.Context(), req.VACode, req.VAName)

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

func GetVAConfigs(cfgSvc *common.VAConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		claims := ctxutil.GetUserClaims(r.Context())
		cfgs, _ := cfgSvc.GetAllConfigValues(r.Context(), claims.ServerID())

		resp := dtos.APIResponse{
			Status:       strconv.FormatBool(true),
			ResponseTime: common.GetResponseTime(initTime),
			Data:         cfgs,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func ListConfigKeys(cfgSvc *common.VAConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()
		cfgs := cfgSvc.ListPossibleKeys()

		resp := dtos.APIResponse{
			Status:       strconv.FormatBool(true),
			ResponseTime: common.GetResponseTime(initTime),
			Data: dtos.VAConfigKeys{
				ConfigKeys: cfgs,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func SetConfigKeys(cfgSvc *common.VAConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		cfgs := make(map[string]string)

		if err := json.NewDecoder(r.Body).Decode(&cfgs); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		res, err := cfgSvc.SetVaConfig(r.Context(), cfgs)

		msg := "Config set successfully"

		if err != nil {
			fmt.Printf("\nPANIC | ERROR \n%v\n", err)
			msg = err.Error()
		}

		resp := dtos.APIResponse{
			Status:       strconv.FormatBool(true),
			ResponseTime: common.GetResponseTime(initTime),
			Message:      msg,
			Data:         res,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func SyncUser(mgmtSvc *services.VAManagementService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		var req dtos.SyncUser
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		status, err := mgmtSvc.SyncUser(r.Context(), req.UserID, req.Callsign)
		if err != nil {
			resp := dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      status,
				ResponseTime: common.GetResponseTime(initTime),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := dtos.APIResponse{
			Status:       status,
			ResponseTime: common.GetResponseTime(initTime),
			Message:      status,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func SetRole(mgmtSvc *services.VAManagementService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		log.Printf("Request received")
		var req dtos.SetRole
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		_, err := mgmtSvc.UpdateUserRole(r.Context(), req.UserID, req.Role)

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
			ResponseTime: common.GetResponseTime(initTime),
			Message:      "Role updated!",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

	}
}
