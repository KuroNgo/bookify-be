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
	ShortDescription  string             `bson:"short_description" json:"short_description"`
	Description       string             `bson:"description" json:"description"`
	ImageURL          string             `bson:"image_url" json:"image_url"`
	AssetURL          string             `bson:"asset_url" json:"asset_url"`
	StartTime         time.Time          `bson:"start_time" json:"start_time"`
	EndTime           time.Time          `bson:"end_time" json:"end_time"`
	Mode              string             `bson:"mode" json:"mode"` // Public, Friend and Group, Invite Only
	EstimatedAttendee int16              `bson:"estimated_attendee" json:"estimated_attendee"`
	ActualAttendee    int16              `bson:"actual_attendee" json:"actual_attendee"`
	TotalExpenditure  float64            `bson:"total_expenditure" json:"total_expenditure"`
	Tags              []string           `bson:"tags" json:"tags"`
}

type EventInput struct {
	EventTypeName     string   `bson:"event_type_name" json:"event_type_name"`
	OrganizationID    string   `bson:"organization_id" json:"organization_id"`
	Title             string   `bson:"title" json:"title"`
	ShortDescription  string   `bson:"short_description" json:"short_description"`
	Description       string   `bson:"description" json:"description"`
	ImageURL          string   `bson:"image_url" json:"image_url"`
	AssetURL          string   `bson:"asset_url" json:"asset_url"`
	StartTime         string   `bson:"start_time" json:"start_time"`
	EndTime           string   `bson:"end_time" json:"end_time"`
	Mode              string   `bson:"mode" json:"mode"` // Public, Friend and Group, Invite Only
	EstimatedAttendee int16    `bson:"estimated_attendee" json:"estimated_attendee"`
	ActualAttendee    int16    `bson:"actual_attendee" json:"actual_attendee"`
	TotalExpenditure  float64  `bson:"total_expenditure" json:"total_expenditure"`
	Tags              []string `bson:"tags" json:"tags"`
	Capacity          int32    `bson:"capacity" json:"capacity"`
	AddressLine       string   `bson:"address_line" json:"address_line"`
	City              string   `bson:"city" json:"city"`
	Country           string   `bson:"country" json:"country"`
	EventMode         string   `bson:"event_mode" json:"event_mode"`
	LinkAttend        string   `bson:"link_attend" json:"link_attend"`
	FromAttend        string   `bson:"from_attend" json:"from_attend"`
}

type EventResponse struct {
	Event Event `json:"event"`
}

type EventResponsePage struct {
	Event       []Event `json:"event"`
	Page        int64   `json:"page"`
	CurrentPage int     `json:"current_page"`
}

type EventResponseForUser struct {
	EventResponse EventResponse `bson:"event_response" json:"event_response"`
	EventType     EventType     `bson:"event_type" json:"event_type"`
	Venue         Venue         `bson:"venue" json:"venue"`
	Organization  Organization  `bson:"organization" json:"organization"`
}
