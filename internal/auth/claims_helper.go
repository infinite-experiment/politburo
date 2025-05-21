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
		return nil
	}
	log.Printf("CreateUser called with discordID=%q username=%q", serverId, userId)
	log.Printf("CreateUser called with discordID=%q username=%q", *user.UserName, user.ID)

	return &APIKeyClaims{
		UserIDValue:   userId,
		ServerIDValue: serverId,
	}
}
