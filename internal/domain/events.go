package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEvent = "event"
)

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	StartTime   time.Time          `bson:"start_time" json:"start_time"`
	EndTime     time.Time          `bson:"end_time" json:"end_time"`
	Location    string             `bson:"location" json:"location"`
	//Capacity    int16              `bson:"capacity" json:"capacity"`
	//Mode        string             `bson:"mode" json:"mode"` // Public, Friend and Group, Invite Only
	AttendanceID primitive.ObjectID `bson:"user_id" json:"user_id"`
	WhoCreated   primitive.ObjectID `bson:"who_created" json:"who_created"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type EventInput struct {
	Name        string `json:"name" bson:"name" example:"Tech Conference 2024" validate:"required"`
	Description string `json:"description" bson:"description" example:"A conference to discuss the latest trends in technology." validate:"required"`
	StartTime   string `json:"start_time" bson:"start_time" example:"2024-05-15 09:00:00" validate:"required,datetime"`
	EndTime     string `json:"end_time" bson:"end_time" example:"2024-05-15 17:00:00" validate:"required,datetime"`
	Location    string `json:"location" bson:"location" example:"San Francisco, CA" validate:"required"`
	UserID      string `json:"user_id" bson:"user_id" example:"6482cfc5f7e3f4b789123456" validate:"required"`
}

type EventResponse struct {
	Event []Event `bson:"event" json:"event"`
}
