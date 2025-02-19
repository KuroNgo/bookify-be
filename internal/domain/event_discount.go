package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEventDiscount = "event_discount"
)

type EventDiscount struct {
	ID              primitive.ObjectID   `bson:"_id" json:"id"`
	EventID         primitive.ObjectID   `bson:"event_id" json:"event_id"`
	DiscountValue   int                  `bson:"discount_value" json:"discount_value"`
	DiscountUnit    string               `bson:"discount_unit" json:"discount_unit"`
	DateCreated     time.Time            `bson:"date_created" json:"date_created"`
	StartDate       time.Time            `bson:"start_date" json:"start_date"`
	ApplicableUsers []primitive.ObjectID `bson:"applicable_users,omitempty" json:"applicable_users,omitempty"` // User được chọn
	EndDate         time.Time            `bson:"end_date" json:"end_date"`
	CreatedAt       time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time            `bson:"updated_at" json:"updated_at"`
	WhoCreated      string               `bson:"who_created" json:"who_created"`
}

type EventDiscountInput struct {
	EventID         string    `bson:"event_id" json:"event_id"`
	DiscountValue   int       `bson:"discount_value" json:"discount_value"`
	DiscountUnit    string    `bson:"discount_unit" json:"discount_unit"`
	DateCreated     time.Time `bson:"date_created" json:"date_created"`
	StartDate       time.Time `bson:"start_date" json:"start_date"`
	ApplicableUsers []string  `bson:"applicable_users,omitempty" json:"applicable_users,omitempty"` // User được chọn
	EndDate         time.Time `bson:"end_date" json:"end_date"`
}
