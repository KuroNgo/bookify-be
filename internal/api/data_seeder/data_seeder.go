package data_seeder

import (
	event_type_data_seeder "bookify/internal/repository/event_type/data_seeder"
	user_seeder "bookify/internal/repository/user/data_seeder"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func DataSeeds(ctx context.Context, client *mongo.Client) error {
	err := user_seeder.SeedUser(ctx, client)
	if err != nil {
		return err
	}

	err = event_type_data_seeder.SeedEventType(ctx, client)
	if err != nil {
		return err
	}

	return nil
}
