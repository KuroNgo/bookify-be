package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEmployeeTicketAssigment = "employee_ticket_assigment"
)

type EventTicketAssignment struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	EventID      primitive.ObjectID `bson:"event_id" json:"event_id"`
	AttendanceID primitive.ObjectID `bson:"attendance_id" json:"attendance_id"`
	PurchaseDate time.Time          `bson:"purchase_date" json:"purchase_date"`
	ExpiryDate   time.Time          `bson:"expiry_date" json:"expiry_date"`
	Price        float64            `bson:"price" json:"price"`
	TicketType   string             `bson:"ticket_type" json:"ticket_type"`
	Status       string             `bson:"status" json:"status"` // use, useless
}

type EventTicketAssignmentInput struct {
	ID           string    `bson:"_id" json:"id"`
	EventID      string    `bson:"event_id" json:"event_id"`
	AttendanceID string    `bson:"attendance_id" json:"attendance_id"`
	PurchaseDate time.Time `bson:"purchase_date" json:"purchase_date"`
	ExpiryDate   time.Time `bson:"expiry_date" json:"expiry_date"`
	Price        float64   `bson:"price" json:"price"`
	TicketType   string    `bson:"ticket_type" json:"ticket_type"`
	Status       string    `bson:"status" json:"status"` // use, useless
}
