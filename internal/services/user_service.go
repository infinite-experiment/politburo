package services

import (
	"context"
	"fmt"
	"log"

	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/models/dtos/responses"
	"infinite-experiment/politburo/internal/models/entities"
)

type UserService struct {
	repo            *repositories.UserRepository
	userRepoGorm    *repositories.UserRepositoryGORM
	pilotStatsService *PilotStatsService
}

func NewUserService(repo *repositories.UserRepository, repoGorm *repositories.UserRepositoryGORM, pilotStatsService *PilotStatsService) *UserService {
	return &UserService{
		repo:            repo,
		userRepoGorm:    repoGorm,
		pilotStatsService: pilotStatsService,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, user *entities.User) error {
	return s.repo.InsertUser(ctx, user)
}

// GetUserDetails retrieves user details with VA affiliations and current VA status
func (s *UserService) GetUserDetails(ctx context.Context, userDiscordID, vaDiscordServerID string) (*responses.UserDetailResponse, error) {
	// Fetch user with all VA affiliations
	user, err := s.userRepoGorm.GetUserWithVAAffiliations(ctx, userDiscordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	// Build affiliations list
	affiliations := make([]responses.VAAffiliation, 0, len(user.UserVARoles))
	var currentVA *responses.CurrentVAStatus

	for _, vaRole := range user.UserVARoles {
		affiliation := responses.VAAffiliation{
			VAID:     vaRole.VAID,
			VAName:   vaRole.VA.Name,
			VACode:   vaRole.VA.Code,
			Role:     string(vaRole.Role),
			IsActive: vaRole.IsActive,
			JoinedAt: vaRole.JoinedAt,
			Callsign: vaRole.Callsign,
		}
		affiliations = append(affiliations, affiliation)

		// Check if this is the current VA (from context)
		// Debug logging
		log.Printf("[GetUserDetails] Comparing VA DiscordID '%s' with context serverID '%s'", vaRole.VA.DiscordID, vaDiscordServerID)

		if vaRole.VA.DiscordID == vaDiscordServerID {
			currentVA = &responses.CurrentVAStatus{
				IsMember: true,
				Role:     string(vaRole.Role),
				IsActive: vaRole.IsActive,
				Callsign: vaRole.Callsign,
			}

			// Fetch Airtable data if user is a member and has a role
			if currentVA.IsMember && currentVA.Role != "" && s.pilotStatsService != nil {
				pilotStatus, err := s.pilotStatsService.GetPilotStatusByCallsign(ctx, userDiscordID, vaRole.VAID)
				if err != nil {
					// Log error but don't fail the request - Airtable data is optional
					log.Printf("[GetUserDetails] Failed to fetch Airtable data for user %s: %v", userDiscordID, err)
				} else {
					// Add Airtable data to current VA status
					currentVA.AirtableData = pilotStatus.RawFields
				}
			}
		}
	}

	// If current VA not found in affiliations, set IsMember to false
	if currentVA == nil {
		currentVA = &responses.CurrentVAStatus{
			IsMember: false,
			IsActive: false,
		}
	}

	// Build response
	response := &responses.UserDetailResponse{
		UserID:        user.ID,
		DiscordID:     user.DiscordID,
		IFCommunityID: user.IFCommunityID,
		IFApiID:       user.IFApiID,
		UserName:      user.UserName,
		IsActive:      user.IsActive,
		CreatedAt:     user.CreatedAt,
		Affiliations:  affiliations,
		CurrentVA:     currentVA,
	}

	return response, nil
}
