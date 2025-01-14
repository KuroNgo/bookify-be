package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Payment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookingID     string             `bson:"booking_id" json:"booking_id"`
	PaymentMethod string             `bson:"payment_method" json:"payment_method"`
	Amount        float64            `bson:"amount" json:"amount"`
	Status        string             `bson:"status" json:"status"` // Pending, Completed, Failed
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
