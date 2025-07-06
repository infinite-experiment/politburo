package constants

const (
	GetUserByDiscordId = `
	SELECT * FROM users WHERE discord_id = $1
	`

	DeleteAllUsers = `
	DELETE FROM users where discord_id IS NOT NULL
	`

	GetStatusByApiKey = `
	SELECT id, status from api_keys where id = $1
	`

	InsertUser = `
	INSERT INTO users (
	discord_id,
	if_community_id,
	if_api_id,
	is_active
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at;
	`

	InsertVA = `
	INSERT INTO virtual_airlines (
		name,
		code,
		discord_server_id,
		callsign_prefix,
		callsign_suffix,
		is_active
	) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at;
	`

	GetVAByDiscordID = `
	SELECT * FROM virtual_airlines where id = $1
	`

	InsertMembership = `
	INSERT INTO va_user_roles (
		user_id,
		va_id,
		role
	) VALUES ($1, $2, $3)
	RETURNING id, joined_at;
	`

	GetUserMembership = `
	WITH u AS (
		SELECT id
		FROM   users
		WHERE  discord_id = $1            -- X-Discord-Id
		LIMIT  1
	),
	v AS (
		SELECT id
		FROM   virtual_airlines
		WHERE  discord_server_id = $2     -- X-Server-Id
		LIMIT  1
	),
	r AS (
		SELECT role
		FROM   va_user_roles
		WHERE  user_id = (SELECT id FROM u)
		AND  va_id   = (SELECT id FROM v)
		LIMIT  1
	)
	SELECT
		(SELECT id   FROM u) AS user_id,   -- NULL if user absent
		(SELECT id   FROM v) AS va_id,     -- NULL if VA   absent
		(SELECT role FROM r) AS role;      -- NULL if no link

	`
)
