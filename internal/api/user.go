package api

import (
	"encoding/json"
	"net/http"
	"time"

	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/models/dtos/requests"
	"infinite-experiment/politburo/internal/models/dtos/responses"
	"infinite-experiment/politburo/internal/models/entities"
	"infinite-experiment/politburo/internal/services"
)

var userService *services.UserService

func SetUserService(service *services.UserService) {
	userService = service
}

// RegisterUserHandler handles POST /api/v1/user/register
//
// @Summary Register a new user
// @Description Creates a new user with Discord ID and IF Community username.
// @Tags Users
// @Accept json
// @Produce json
// @Param input body requests.RegisterUserRequest true "User Registration Payload"
// @Success 200 {object} dtos.UserRegisterSwaggerResponse
// @Failure 400 {object} dtos.UserRegisterSwaggerResponse
// @Router /api/v1/user/register [post]
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.DiscordID == "" || req.IFCommunityID == "" || req.ServerID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	ctx := r.Context()

	user := &entities.User{
		DiscordID:     req.DiscordID,
		IFCommunityID: req.IFCommunityID,
		IsActive:      false,
	}

	if err := userService.RegisterUser(ctx, user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not register user")
		return
	}

	resp := responses.UserRegisterResponse{
		RegistrationStatus: true,
	}

	respondWithSuccess(w, http.StatusOK, &resp)
}

// InitUserRegistrationHandler handles POST /api/v1/user/register/init
//
// @Summary      Initiate user registration
// @Description  Initiates the user registration process given an IF Community ID (IFC ID).
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param X-Discord-Id header string true "Discord ID"
// @Param X-Server-Id header string true "Discord Server ID"  default(123456789)
// @Param X-API-Key header string true "API KEY"
// @Param        input  body      dtos.InitUserRegistrationReq  true  "IFC ID Payload"
// @Success      200    {object}  dtos.APIResponse
// @Failure      400    {object}  dtos.APIResponse
// @Router       /api/v1/user/register/init [post]
func InitUserRegistrationHandler(regService *services.RegistrationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		initTime := time.Now()
		var req dtos.InitUserRegistrationReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IfcId == "" {

			resp := dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      "Invalid IFC ID Received",
				ResponseTime: common.GetResponseTime(initTime),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		apiResp, _, err := regService.InitUserRegistration(r.Context(), req.IfcId, req.LastFlight)
		if err != nil {
			resp := dtos.APIResponse{
				Status:       string(constants.APIStatusError),
				Message:      "Failed to process",
				ResponseTime: common.GetResponseTime(initTime),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := dtos.APIResponse{
			Status:       string(constants.APIStatusOk),
			Message:      "Initiated",
			ResponseTime: common.GetResponseTime(initTime),
			Data:         apiResp,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

	}

}

// DeleteAllUsers godoc
// @Summary      Delete all users (Test Only)
// @Description  Deletes all users in the database. Intended for development/testing only.
// @Tags         Test
// @Param X-Discord-Id header string true "Discord ID" default(668664447950127154)
// @Param X-Server-Id header string true "Discord Server ID" default(988020008665882624)
// @Param X-API-Key header string true "API KEY" default(API_KEY_123)
// @Produce      json
// @Success      400  {object}  dtos.APIResponse  "Always returns error; not implemented for production use"
// @Router       /api/v1/users/delete [get]
func DeleteAllUsers(repo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo.DeleteAllUsers(r.Context())
		resp := dtos.APIResponse{
			Status:  string(constants.APIStatusError),
			Message: "Invalid IFC ID Received",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)

	}
}
