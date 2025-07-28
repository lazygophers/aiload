package model

import "gorm.io/gorm"

// User represents the user model.
type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}
