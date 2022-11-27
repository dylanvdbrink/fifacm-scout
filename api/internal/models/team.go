package models

import "gorm.io/gorm"

type Team struct {
	gorm.Model `json:"-"`

	TeamID  int `gorm:"unique;"`
	Name    string
	LogoUrl string

	DBUpdateID uint     `json:"-"`
	DBUpdate   DBUpdate `json:"-"`
}
