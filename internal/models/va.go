package models

import "time"

type VA struct {
	ID        string    `gorm:"column:id;primaryKey"` // mark primary key
	DiscordID string    `gorm:"column:discord_server_id"`
	Name      string    `gorm:"column:name"`
	Code      string    `gorm:"column:code"`
	IsActive  bool      `gorm:"column:is_active"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// If your table is named differently (not "v_a" or "vas"), add TableName() method
func (VA) TableName() string {
	return "virtual_airlines"
}
