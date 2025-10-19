package responses

import "time"

// UserDetailResponse represents detailed user information including VA affiliations
type UserDetailResponse struct {
	UserID        string          `json:"user_id"`
	DiscordID     string          `json:"discord_id"`
	IFCommunityID string          `json:"if_community_id"`
	IFApiID       *string         `json:"if_api_id,omitempty"`
	UserName      *string         `json:"username,omitempty"`
	IsActive      bool            `json:"is_active"`
	CreatedAt     time.Time       `json:"created_at"`
	Affiliations  []VAAffiliation `json:"affiliations"`
	CurrentVA     *CurrentVAStatus `json:"current_va,omitempty"`
}

// VAAffiliation represents a user's membership in a virtual airline
type VAAffiliation struct {
	VAID     string    `json:"va_id"`
	VAName   string    `json:"va_name"`
	VACode   string    `json:"va_code"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
	JoinedAt time.Time `json:"joined_at"`
	Callsign string    `json:"callsign,omitempty"`
}

// CurrentVAStatus represents the user's status in the current VA context
type CurrentVAStatus struct {
	IsMember bool   `json:"is_member"`
	Role     string `json:"role,omitempty"`
	IsActive bool   `json:"is_active"`
	Callsign string `json:"callsign,omitempty"`
}
