package auth

import (
	"context"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/db/repositories"
	"log"
)

func MakeClaimsFromApi(ctx context.Context, userRepo *repositories.UserRepository, serverId string, userId string) *APIKeyClaims {
	log.Printf("Checking for user: %q", userId)

	member, err := userRepo.FindUserMembership(ctx, userId, serverId)
	if err != nil {
		log.Printf("membership lookup error: %v", err)
		// Return a minimal claims object; UUIDs stay empty
		return &APIKeyClaims{
			DiscordUIDVal:      userId,
			DiscordServerIDVal: serverId,
		}
	}

	log.Printf("\nMembership result\n%v\n", member)

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

	log.Printf("User ID: %s, VA ID: %s, Role: %s, U UUID: %s, VA UUID: %s", userId, serverId, role, userUUID, vaUUID)

	return &APIKeyClaims{
		UserUUID:           userUUID,
		VaUUID:             vaUUID,
		RoleValue:          role,
		DiscordUIDVal:      userId,
		DiscordServerIDVal: serverId,
	}
}
