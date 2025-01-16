package partner_repository

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

type IPartnerRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.Partner, error)
	GetAll(ctx context.Context) ([]domain.Partner, error)
	CreateOne(ctx context.Context, partner *domain.Partner) error
	UpdateOne(ctx context.Context, partner *domain.Partner) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	CountExist(ctx context.Context, name string) (int64, error)
}

type partnerRepository struct {
	database          *mongo.Database
	collectionPartner string
}

func (p partnerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.Partner, error) {
	partnerCollection := p.database.Collection(p.collectionPartner)

	filter := bson.M{"_id": id}
	var partner domain.Partner
	if err := partnerCollection.FindOne(ctx, filter).Decode(&partner); err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return domain.Partner{}, nil
		}
		return domain.Partner{}, err
	}

	return partner, nil
}

func (p partnerRepository) GetAll(ctx context.Context) ([]domain.Partner, error) {
	partnerCollection := p.database.Collection(p.collectionPartner)

	filter := bson.M{}
	cursor, err := partnerCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var partners []domain.Partner
	for cursor.Next(ctx) {
		var partner domain.Partner
		if err := cursor.Decode(&partner); err != nil {
			return nil, err
		}

		partners = append(partners, partner)
	}

	return partners, nil
}

func (p partnerRepository) CreateOne(ctx context.Context, partner *domain.Partner) error {
	partnerCollection := p.database.Collection(p.collectionPartner)

	if err := validate_data.ValidatePartner(partner); err != nil {
		return err
	}

	_, err := partnerCollection.InsertOne(ctx, partner)
	if err != nil {
		return err
	}

	return nil
}

func (p partnerRepository) UpdateOne(ctx context.Context, partner *domain.Partner) error {
	partnerCollection := p.database.Collection(p.collectionPartner)

	if err := validate_data.ValidatePartner(partner); err != nil {
		return err
	}

	filter := bson.M{"_id": partner.ID}
	update := bson.M{"$set": bson.M{
		"name":  partner.Name,
		"phone": partner.Phone,
		"email": partner.Email,
	}}
	_, err := partnerCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (p partnerRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	partnerCollection := p.database.Collection(p.collectionPartner)

	if id == primitive.NilObjectID {
		return errors.New(constants.MsgDataInvalidFormat)
	}

	filter := bson.M{"_id": id}
	_, err := partnerCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (p partnerRepository) CountExist(ctx context.Context, name string) (int64, error) {
	partnerCollection := p.database.Collection(p.collectionPartner)

	filter := bson.M{"name": name}
	count, err := partnerCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func NewPartnerRepository(database *mongo.Database, collectionPartner string) IPartnerRepository {
	return &partnerRepository{database: database, collectionPartner: collectionPartner}
}
