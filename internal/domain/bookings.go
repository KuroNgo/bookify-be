package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	EventID   primitive.ObjectID `bson:"event_id" json:"event_id"`
	TicketID  primitive.ObjectID `bson:"ticket_id" json:"ticket_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Status    string             `bson:"status" json:"status"` // Pending, Confirmed, Cancelled
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type BookingInput struct {
	UserID   string `bson:"user_id" json:"user_id"`
	EventID  string `bson:"event_id" json:"event_id"`
	TicketID string `bson:"ticket_id" json:"ticket_id"`
	Quantity int    `bson:"quantity" json:"quantity"`
	Status   string `bson:"status" json:"status"` // Pending, Confirmed, Cancelled
}

type BookingResponse struct {
	Booking []Booking `bson:"booking" json:"booking"`
}
