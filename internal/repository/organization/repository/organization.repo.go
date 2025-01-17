package organization_repository

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/constants"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IOrganizationRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.Organization, error)
	GetAll(ctx context.Context) ([]domain.Organization, error)
	CreateOne(ctx context.Context, partner *domain.Organization) error
	UpdateOne(ctx context.Context, partner *domain.Organization) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	CountExist(ctx context.Context, name string) (int64, error)
}

type organizationRepository struct {
	database               *mongo.Database
	collectionOrganization string
}

func NeOrganizationRepository(database *mongo.Database, collectionOrganization string) IOrganizationRepository {
	return &organizationRepository{database: database, collectionOrganization: collectionOrganization}
}

func (o organizationRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.Organization, error) {
	organizationCollection := o.database.Collection(o.collectionOrganization)

	filter := bson.M{"_id": id}
	var organization domain.Organization
	if err := organizationCollection.FindOne(ctx, filter).Decode(&organization); err != nil {
		return domain.Organization{}, err
	}

	return organization, nil
}

func (o organizationRepository) GetAll(ctx context.Context) ([]domain.Organization, error) {
	organizationCollection := o.database.Collection(o.collectionOrganization)

	filter := bson.M{}
	cursor, err := organizationCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var organizations []domain.Organization
	for cursor.Next(ctx) {
		var organization domain.Organization
		if err = cursor.Decode(&organization); err != nil {
			return nil, err
		}

		organizations = append(organizations, organization)
	}

	return organizations, nil
}

func (o organizationRepository) CreateOne(ctx context.Context, organization *domain.Organization) error {
	organizationCollection := o.database.Collection(o.collectionOrganization)

	if err := validate_data.ValidateOrganization(organization); err != nil {
		return err
	}

	filter := bson.M{"name": organization.Name}
	count, err := organizationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	_, err = organizationCollection.InsertOne(ctx, organization)
	if err != nil {
		return err
	}

	return nil
}

func (o organizationRepository) UpdateOne(ctx context.Context, organization *domain.Organization) error {
	organizationCollection := o.database.Collection(o.collectionOrganization)

	if err := validate_data.ValidateOrganization(organization); err != nil {
		return err
	}

	filter := bson.M{"name": organization.Name}
	update := bson.M{"$set": bson.M{
		"name":           organization.Name,
		"contact_person": organization.ContactPerson,
		"email":          organization.Email,
		"phone":          organization.Phone,
	}}
	count, err := organizationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(constants.MsgAPIConflict)
	}

	_, err = organizationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (o organizationRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	organizationCollection := o.database.Collection(o.collectionOrganization)

	filter := bson.M{"_id": id}
	_, err := organizationCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (o organizationRepository) CountExist(ctx context.Context, name string) (int64, error) {
	organizationCollection := o.database.Collection(o.collectionOrganization)

	filter := bson.M{"name": name}
	count, err := organizationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
