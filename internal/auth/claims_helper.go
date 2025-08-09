package auth

import (
	"context"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/db/repositories"
)

func MakeClaimsFromApi(ctx context.Context, userRepo *repositories.UserRepository, serverId string, userId string) *APIKeyClaims {

	member, err := userRepo.FindUserMembership(ctx, serverId, userId)
	if err != nil {
		// Return a minimal claims object; UUIDs stay empty
		return &APIKeyClaims{
			DiscordUIDVal:      userId,
			DiscordServerIDVal: serverId,
		}
	}

	if member == nil { // no row found
		return &APIKeyClaims{
			DiscordUIDVal:      userId,
			DiscordServerIDVal: serverId,
		}
	}

	var userUUID, vaUUID string
	var role constants.VARole
	if member.UserID != nil {
		userUUID = *member.UserID
	}
	if member.VAID != nil {
		vaUUID = *member.VAID
	}
	if member.Role != nil {
		role = *member.Role
	}

	return &APIKeyClaims{
		UserUUID:           userUUID,
		VaUUID:             vaUUID,
		RoleValue:          role,
		DiscordUIDVal:      userId,
		DiscordServerIDVal: serverId,
	}
}
