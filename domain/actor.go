package domain

import "time"

type Actor struct {
	ID         int64     `json:"ID" gorm:"primaryKey;type:SERIAL;"`
	FirstName  string    `json:"FirstName" gorm:"type:VARCHAR(100);NOT NULL;"`
	LastName   string    `json:"LastName" gorm:"type:VARCHAR(100);NOT NULL;"`
	LastUpdate time.Time `json:"LastUpdate" gorm:"default:CURRENT_TIMESTAMP;NOT NULL;"`
}
