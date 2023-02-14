package entities

import (
	"gorm.io/gorm"
)

type User struct {
	//ID        uint   `json:"uid"`
	gorm.Model
	Name       string `json:"name"`
	Email      string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password   string `gorm:"not null" json:"password" binding:"required"`
	Address    string `json:"addr"`
	DOB        string `json:"dob"`
	Role string `json:"role"`
}
