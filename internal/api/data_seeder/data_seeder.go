package data_seeder

import (
	eventtypedataseeder "bookify/internal/repository/event_type/data_seeder"
	organizationseeder "bookify/internal/repository/organization/data_seeder"
	userseeder "bookify/internal/repository/user/data_seeder"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func DataSeeds(ctx context.Context, client *mongo.Client) error {
	err, userID := userseeder.SeedUser(ctx, client)
	if err != nil {
		return err
	}

	err = eventtypedataseeder.SeedEventType(ctx, client)
	if err != nil {
		return err
	}

	err = organizationseeder.SeedOrganization(ctx, client, userID)
	if err != nil {
		return err
	}

	return nil
}
