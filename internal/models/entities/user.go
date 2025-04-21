package entities

import "time"

type User struct {
	ID            string    `db:"id"`
	DiscordID     string    `db:"discord_id"`
	IFCommunityID string    `db:"if_community_id"`
	IFApiID       *string   `db:"if_api_id"`
	IsActive      bool      `db:"is_active"`
	UserName      *string   `db:"username"`
	OTP           *string   `db:"otp"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
