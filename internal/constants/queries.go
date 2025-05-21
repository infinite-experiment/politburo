package constants

const (
	GetUserByDiscordId = `
	SELECT * FROM users WHERE discord_id = $1
	`
)
