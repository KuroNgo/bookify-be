package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionEvent = "event"
)

type Event struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventTypeID       primitive.ObjectID `bson:"event_type_id" json:"event_type_id"`
	VenueID           primitive.ObjectID `bson:"venue_id" json:"venue_id"`
	OrganizationID    primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	Title             string             `bson:"title" json:"title"`
	Description       string             `bson:"description" json:"description"`
	ImageURL          string             `bson:"image_url" json:"image_url"`
	AssetURL          string             `bson:"asset_url" json:"asset_url"`
	StartTime         time.Time          `bson:"start_time" json:"start_time"`
	EndTime           time.Time          `bson:"end_time" json:"end_time"`
	Mode              string             `bson:"mode" json:"mode"` // Public, Friend and Group, Invite Only
	EstimatedAttendee int16              `bson:"estimated_attendee" json:"estimated_attendee"`
	ActualAttendee    int16              `bson:"actual_attendee" json:"actual_attendee"`
	TotalExpenditure  float64            `bson:"total_expenditure" json:"total_expenditure"`
}

type EventInput struct {
	EventTypeID       primitive.ObjectID `bson:"event_type_id" json:"event_type_id"`
	VenueID           primitive.ObjectID `bson:"venue_id" json:"venue_id"`
	OrganizationID    primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	Title             string             `bson:"title" json:"title"`
	Description       string             `bson:"description" json:"description"`
	ImageURL          string             `bson:"image_url" json:"image_url"`
	AssetURL          string             `bson:"asset_url" json:"asset_url"`
	StartTime         time.Time          `bson:"start_time" json:"start_time"`
	EndTime           time.Time          `bson:"end_time" json:"end_time"`
	Mode              string             `bson:"mode" json:"mode"` // Public, Friend and Group, Invite Only
	EstimatedAttendee int16              `bson:"estimated_attendee" json:"estimated_attendee"`
	ActualAttendee    int16              `bson:"actual_attendee" json:"actual_attendee"`
	TotalExpenditure  float64            `bson:"total_expenditure" json:"total_expenditure"`
}

type EventResponse struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventTypeID    primitive.ObjectID `bson:"event_type_id" json:"event_type_id"`
	VenueID        primitive.ObjectID `bson:"venue_id" json:"venue_id"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	Title          string             `bson:"title" json:"title"`
	Description    string             `bson:"description" json:"description"`
	ImageURL       string             `bson:"image_url" json:"image_url"`
	AssetURL       string             `bson:"asset_url" json:"asset_url"`
	StartTime      time.Time          `bson:"start_time" json:"start_time"`
	EndTime        time.Time          `bson:"end_time" json:"end_time"`
	Mode           string             `bson:"mode" json:"mode"` // Public, Friend and Group, Invite Only
}

type EventResponseForUser struct {
	EventResponse EventResponse `bson:"event_response" json:"event_response"`
	EventType     EventType     `bson:"event_type" json:"event_type"`
	Venue         Venue         `bson:"venue" json:"venue"`
	Organization  Organization  `bson:"organization" json:"organization"`
}
