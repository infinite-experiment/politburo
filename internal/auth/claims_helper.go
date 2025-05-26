package auth

import (
	"context"
	"infinite-experiment/politburo/internal/db/repositories"
	"log"
)

func MakeClaimsFromApi(ctx context.Context, repo *repositories.UserRepository, serverId, userId string) *APIKeyClaims {

	log.Printf("Checking for user: %q", userId)
	user, err := repo.FindUserByDiscordId(ctx, userId)
	if err != nil {
		log.Print(err.Error())
	}
	log.Printf("CreateUser called with discordID=%q ", serverId)

	log.Printf("DB Query response: %v", user)
	return &APIKeyClaims{
		UserIDValue:   userId,
		ServerIDValue: serverId,
	}
}
