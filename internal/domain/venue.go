package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionVenue = "venue"
)

type Venue struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Capacity    int32              `bson:"capacity" json:"capacity"`
	AddressLine string             `bson:"address_line" json:"address_line"`
	City        string             `bson:"city" json:"city"`
	Country     string             `bson:"country" json:"country"`
	EventMode   string             `bson:"event_mode" json:"event_mode"`
	LinkAttend  string             `bson:"link_attend" json:"link_attend"`
	FromAttend  string             `bson:"from_attend" json:"from_attend"`
}

type VenueInput struct {
	Capacity    int32  `bson:"capacity" json:"capacity"`
	AddressLine string `bson:"address_line" json:"address_line"`
	City        string `bson:"city" json:"city"`
	Country     string `bson:"country" json:"country"`
	EventMode   string `bson:"event_mode" json:"event_mode"`
	LinkAttend  string `bson:"link_attend" json:"link_attend"`
	FromAttend  string `bson:"from_attend" json:"from_attend"`
}
