package domain

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Message string
	UserId  string
	Status  string
}
type AppointmentEvent struct {
	AppointmentId   int
	Email string
	VideoURL string
	DoctorId        string
	AppointmentDate string
	Type            string
}