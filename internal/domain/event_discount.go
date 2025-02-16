package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEventDiscount = "event_discount"
)

type EventDiscount struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	EventID       primitive.ObjectID `bson:"event_id" json:"event_id"`
	DiscountValue int                `bson:"discount_value" json:"discount_value"`
	DiscountUnit  string             `bson:"discount_unit" json:"discount_unit"`
	DateCreated   time.Time          `bson:"date_created" json:"date_created"`
}

type EventDiscountInput struct {
	EventID       primitive.ObjectID `bson:"event_id" json:"event_id"`
	DiscountValue int                `bson:"discount_value" json:"discount_value"`
	DiscountUnit  string             `bson:"discount_unit" json:"discount_unit"`
	DateCreated   time.Time          `bson:"date_created" json:"date_created"`
}
