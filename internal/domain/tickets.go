package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Ticket struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventID    string             `bson:"event_id" json:"event_id"`
	TicketType string             `bson:"ticket_type" json:"ticket_type"`
	Price      float64            `bson:"price" json:"price"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
