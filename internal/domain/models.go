package domain

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Message string
	UserId  string
	Status  string
}
