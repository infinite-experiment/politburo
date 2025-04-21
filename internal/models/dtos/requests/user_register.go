package requests

type RegisterUserRequest struct {
	DiscordID     string `json:"discord_id" validate:"required"`
	IFCommunityID string `json:"if_community_id" validate:"required"`
	ServerID      string `json:"server_id" validate:"required"`
}
