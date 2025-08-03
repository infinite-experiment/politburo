package constants

const (
	GetUserByDiscordId = `
	SELECT * FROM users WHERE discord_id = $1
	`

	DeleteAllUsers = `
	DELETE FROM users where discord_id IS NOT NULL
	`

	DeleteAllRoles = `
	DELETE FROM va_user_roles where id IS NOT NULL
	`

	DeleteAllServers = `
	DELETE FROM virtual_airlines where id IS NOT NULL
	`

	GetVAConfigs = `
	SELECT id, va_id, config_key, config_value, created_at, updated_at
	FROM va_configs
	WHERE va_id = $1	`

	UpsertVAConfig = `
	INSERT INTO va_configs (va_id, config_key, config_value)
	VALUES ($1, $2, $3)
	ON CONFLICT (va_id, config_key)
	DO UPDATE SET
		config_value = EXCLUDED.config_value,
		updated_at = NOW();
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
		is_active
	) VALUES ($1, $2, $3, $4)
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
