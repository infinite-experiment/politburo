package api

import (
	"encoding/json"
	"net/http"

	"infinite-experiment/infinite-experiment-backend/internal/models/dtos/requests"
	"infinite-experiment/infinite-experiment-backend/internal/models/dtos/responses"
	"infinite-experiment/infinite-experiment-backend/internal/models/entities"
	"infinite-experiment/infinite-experiment-backend/internal/services"
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
