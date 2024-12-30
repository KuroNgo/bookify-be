package user_seeder

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/password"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var user = domain.User{
	ID:           primitive.NewObjectID(),
	Username:     "admin",
	Email:        "admin@admin.com",
	PasswordHash: "12345",
	Phone:        "0329245971",
	Verified:     true,
	Provider:     "app",
	Role:         constants.RoleSuperAdmin,
	CreatedAt:    time.Now(),
	UpdatedAt:    time.Now(),
}

func SeedUser(ctx context.Context, client *mongo.Client) error {
	collectionUser := client.Database("bookify").Collection("user")

	count, err := collectionUser.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	user.PasswordHash, err = password.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = collectionUser.InsertOne(ctx, user)
		if err != nil {
			return err
		}
	}

	return nil
}
