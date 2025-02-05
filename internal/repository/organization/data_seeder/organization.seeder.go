package data_seeder

import (
	"bookify/internal/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedOrganization(ctx context.Context, client *mongo.Client, userID primitive.ObjectID) error {
	collectionOrganization := client.Database("bookify").Collection("organization")

	var organizationInput = domain.Organization{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		Name:          "Tech Innovators",
		ContactPerson: "John Doe",
		Email:         "johndoe@example.com",
		Phone:         "+1234567890",
	}

	count, err := collectionOrganization.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = collectionOrganization.InsertOne(ctx, organizationInput)
		if err != nil {
			return err
		}
	}

	return nil
}
