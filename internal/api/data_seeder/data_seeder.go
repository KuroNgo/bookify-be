package data_seeder

import (
	user_seeder "bookify/internal/repository/user/data_seeder"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func DataSeeds(ctx context.Context, client *mongo.Client) error {
	err := user_seeder.SeedUser(ctx, client)
	if err != nil {
		return err
	}

	return nil
}
