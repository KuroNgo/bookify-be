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
	FullName:     "admin",
	Gender:       "male",
	Vocation:     "Marketing",
	Address:      "Thu Duc",
	City:         "Ho Chi Minh city",
	Region:       "Viet Nam",
	DateOfBirth:  time.Date(2002, 2, 6, 16, 5, 0, 0, time.UTC), // Sử dụng time.Time
	Email:        "admin@admin.com",
	PasswordHash: "12345",
	Phone:        "0329245971",
	Verified:     true,
	Provider:     "app",
	Role:         constants.RoleSuperAdmin,
	ShowInterest: false,
	SocialMedia:  false,
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
