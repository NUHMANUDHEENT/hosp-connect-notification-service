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

type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	PatientID string  `json:"patient_id"`
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Date      string  `json:"date"`
}
