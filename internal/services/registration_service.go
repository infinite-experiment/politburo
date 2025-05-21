package services

import (
	stdContext "context"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/models/entities"
)

type RegistrationService struct {
	Cache          CacheService
	LiveAPI        *LiveAPIService
	UserRepository repositories.UserRepository
}

func NewRegistrationService(liveAPI *LiveAPIService, cache CacheService, userRepo repositories.UserRepository) *RegistrationService {
	return &RegistrationService{
		LiveAPI:        liveAPI,
		Cache:          cache,
		UserRepository: userRepo,
	}
}

type InitRegistrationValidation struct {
	UserDB    entities.User
	IFProfile dtos.UserStatsResponse
}

func (svc *RegistrationService) InitUserRegistration(ctx stdContext.Context, ifcId string) (*dtos.InitApiResponse, string, error) {
	// claims := customContext.GetUserClaims(ctx)

	return &dtos.InitApiResponse{
		IfcId:                   ifcId,
		IsVerificationInitiated: true,
		Message:                 "Initiation has started",
		LastFlight:              "AAA",
	}, "", nil

}
