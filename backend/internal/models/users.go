package models

import (
	_ "gorm.io/gorm"
)

type userType string

const (
	PrimaryOwner userType = "primary_owner"
	Client       userType = "client"
)

// User model info
// @Description User model to create a new user
type User struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" swaggerignore:"true"`
	Name           string    `gorm:"size:50;not null" validate:"required"`
	Email          string    `gorm:"size:100;unique;not null" validate:"required,email"`
	Password       string    `gorm:"-" validate:"required"`
	PasswordHash   string    `gorm:"size:60;not null" swaggerignore:"true"`
	UserType       userType  `gorm:"type:user_type;not null" validate:"required"`
	CompanyName    string    `gorm:"size:100" validate:"required"`
	PrimaryOwnerID uint      `validate:"required"`
	Addresses      []Address `gorm:"type:json" swaggerignore:"true"` // Address struct for JSON serialization
	CreatedAt      int       `gorm:"autoCreateTime"`
	UpdatedAt      int       `gorm:"autoUpdateTime"`
}

// func (ct *userType) Scan(value interface{}) error {
// 	*ct = userType(value.([]byte))
// 	return nil
// }

// func (ct userType) Value() (driver.Value, error) {
// 	return string(ct), nil
// }

type Address struct {
	Name         string
	AddressLine1 string `validate:"required"`
	AddressLine2 string
	City         string `validate:"required"`
	State        string `validate:"required"`
	PostalCode   int    `validate:"required"`
	Phone        string `validate:"required"`
}
