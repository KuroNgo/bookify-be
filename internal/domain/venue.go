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
	State       string             `bson:"state" json:"state"`
	Country     string             `bson:"country" json:"country"`
	PostalCode  string             `bson:"postal_code" json:"postal_code"`
	OnlineFlat  bool               `bson:"online_flat" json:"online_flat"`
}

type VenueInput struct {
	Capacity    int32  `bson:"capacity" json:"capacity"`
	AddressLine string `bson:"address_line" json:"address_line"`
	City        string `bson:"city" json:"city"`
	State       string `bson:"state" json:"state"`
	Country     string `bson:"country" json:"country"`
	PostalCode  string `bson:"postal_code" json:"postal_code"`
	OnlineFlat  bool   `bson:"online_flat" json:"online_flat"`
}
