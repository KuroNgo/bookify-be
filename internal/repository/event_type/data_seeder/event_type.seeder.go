package event_type_data_seeder

import (
	"bookify/internal/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var eventTypes = []domain.EventType{
	{
		ID:   primitive.NewObjectID(),
		Name: "Workshop",
	},
	{
		ID:   primitive.NewObjectID(),
		Name: "Webinar",
	},
	{
		ID:   primitive.NewObjectID(),
		Name: "Music Festival",
	},
	{
		ID:   primitive.NewObjectID(),
		Name: "Conference",
	},
}

func SeedEventType(ctx context.Context, client *mongo.Client) error {
	collection := client.Database("bookify").Collection("event_type")

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count == 0 {
		var eventDocs []interface{}
		for _, eventType := range eventTypes {
			eventDocs = append(eventDocs, eventType)
		}
		_, err = collection.InsertMany(ctx, eventDocs)
		if err != nil {
			return err
		}
	}

	return nil
}
