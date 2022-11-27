package models

import "gorm.io/gorm"

type DBUpdate struct {
	gorm.Model `json:"-"`

	UpdateID string
	Name     string
}
