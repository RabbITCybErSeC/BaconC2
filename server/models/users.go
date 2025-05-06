package models

import "gorm.io/gorm"

// User represents a user for authentication
type User struct {
	gorm.Model
	ID       string
	UserName string `gorm:"uniqueIndex"`
	Password string
}
