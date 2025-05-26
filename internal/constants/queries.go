package constants

const (
	GetUserByDiscordId = `
	SELECT * FROM users WHERE discord_id = $1
	`

	DeleteAllUsers = `
	DELETE FROM users where discord_id IS NOT NULL
	`
)
