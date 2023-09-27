package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	Email          string `gorm:"size:100;unique;not null"`
	PasswordHash   string `gorm:"size:60;not null"`
	UserType       string `gorm:"type:enum('primary_owner', 'client');not null"`
	CompanyName    string `gorm:"size:100"`
	CreatedAt      time.Time
	PrimaryOwnerID uint
	Addresses      []Address `gorm:"type:json"` // Address struct for JSON serialization
}

type Address struct {
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   int
}
