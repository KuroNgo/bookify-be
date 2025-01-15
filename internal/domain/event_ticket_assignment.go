package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEmployeeTicketAssigment = "employee_ticket_assigment"
)

type EventTicketAssignment struct {
	TicketID     primitive.ObjectID `bson:"_id" json:"id"`
	EventID      primitive.ObjectID `bson:"event_id" json:"event_id"`
	AttendanceID primitive.ObjectID `bson:"attendance_id" json:"attendance_id"`
	PurchaseDate time.Time          `bson:"purchase_date" json:"purchase_date"`
	ExpiryDate   time.Time          `bson:"expiry_date" json:"expiry_date"`
	Price        int64              `bson:"price" json:"price"`
	TicketType   string             `bson:"ticket_type" json:"ticket_type"`
}
