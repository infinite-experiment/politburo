package api

import (
	"encoding/json"
	"net/http"
	"time"

	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/services"
)

// GetUserDetailsHandler handles GET /api/v1/user/details
//
// @Summary      Get user details with VA affiliations
// @Description  Returns detailed user information including all VA affiliations and current VA status
// @Tags         Users
// @Produce      json
// @Param        X-Discord-Id  header  string  true  "Discord ID"
// @Param        X-Server-Id   header  string  true  "Discord Server ID"
// @Param        X-API-Key     header  string  true  "API KEY"
// @Success      200  {object}  dtos.APIResponse
// @Failure      400  {object}  dtos.APIResponse
// @Failure      500  {object}  dtos.APIResponse
// @Router       /api/v1/user/details [get]
func GetUserDetailsHandler(deps *Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		// Get claims from context
		claims := auth.GetUserClaims(r.Context())
		if claims == nil {
			common.RespondError(w, initTime, nil, "Unauthorized: missing claims", http.StatusUnauthorized)
			return
		}

		userDiscordID := claims.DiscordUserID()
		vaDiscordServerID := claims.DiscordServerID()

		// Call service to get user details
		userDetails, err := deps.Services.User.GetUserDetails(r.Context(), userDiscordID, vaDiscordServerID)
		if err != nil {
			common.RespondError(w, initTime, err, "Failed to fetch user details", http.StatusInternalServerError)
			return
		}

		common.RespondSuccess(w, initTime, "User details fetched successfully", userDetails)
	}
}

// InitUserRegistrationHandler handles POST /api/v1/user/register/init
//
// @Summary      Initiate user registration
// @Description  Initiates the user registration process given an IF Community ID (IFC ID).
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        X-Discord-Id  header  string                        true  "Discord ID"
// @Param        X-Server-Id   header  string                        true  "Discord Server ID"  default(123456789)
// @Param        X-API-Key     header  string                        true  "API KEY"
// @Param        input         body    dtos.InitUserRegistrationReq  true  "IFC ID Payload"
// @Success      200  {object}  dtos.APIResponse
// @Failure      400  {object}  dtos.APIResponse
// @Router       /api/v1/user/register/init [post]
func InitUserRegistrationHandler(regService *services.RegistrationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()
		var req dtos.InitUserRegistrationReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IfcId == "" {
			common.RespondError(w, initTime, err, "Invalid IFC ID Received", http.StatusBadRequest)
			return
		}

		apiResp, _, err := regService.InitUserRegistration(r.Context(), req.IfcId, req.LastFlight)
		if err != nil {
			common.RespondError(w, initTime, err, "Failed to process", http.StatusBadRequest)
			return
		}

		common.RespondSuccess(w, initTime, "Initiated", apiResp)
	}
}

// DeleteAllUsers godoc
// @Summary      Delete all users (Test Only)
// @Description  Deletes all users in the database. Intended for development/testing only.
// @Tags         Test
// @Param        X-Discord-Id  header  string  true  "Discord ID"         default(668664447950127154)
// @Param        X-Server-Id   header  string  true  "Discord Server ID"  default(988020008665882624)
// @Param        X-API-Key     header  string  true  "API KEY"            default(API_KEY_123)
// @Produce      json
// @Success      400  {object}  dtos.APIResponse  "Always returns error; not implemented for production use"
// @Router       /api/v1/users/delete [delete]
func DeleteAllUsers(repo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		repo.DeleteAllUsers(r.Context())
		common.RespondError(w, initTime, nil, "All users deleted", http.StatusBadRequest)
	}
}

// ============================================================================
// Handler Methods (Wrapped for DI pattern - Hybrid Approach)
// ============================================================================

func (h *Handlers) GetUserDetails() http.HandlerFunc {
	return GetUserDetailsHandler(h.deps)
}

func (h *Handlers) InitUserRegistration() http.HandlerFunc {
	return InitUserRegistrationHandler(&h.deps.Services.Reg)
}

func (h *Handlers) DeleteAllUsers() http.HandlerFunc {
	return DeleteAllUsers(&h.deps.Repo.User)
}
