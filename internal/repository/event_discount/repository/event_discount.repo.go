package event_discount_repository

import (
	"bookify/internal/domain"
	"bookify/pkg/shared/validate_data"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IEventDiscountRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventDiscount, error)
	GetByUserIDInApplicableAndExpiringOneDayLeft(ctx context.Context, userID primitive.ObjectID) (domain.EventDiscount, error)
	GetByUserIDInApplicable(ctx context.Context, userID primitive.ObjectID) ([]domain.EventDiscount, error)
	GetByUserIDInApplicableAndEventID(ctx context.Context, userID primitive.ObjectID, eventID primitive.ObjectID) (domain.EventDiscount, error)
	GetAll(ctx context.Context) ([]domain.EventDiscount, error)
	CreateOne(ctx context.Context, discount domain.EventDiscount) error
	UpdateOne(ctx context.Context, discount domain.EventDiscount) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
}

type eventSaleOffRepository struct {
	database                *mongo.Database
	collectionEventDiscount string
}

func NewEventDiscountRepository(database *mongo.Database, collectionEventDiscount string) IEventDiscountRepository {
	return &eventSaleOffRepository{database: database, collectionEventDiscount: collectionEventDiscount}
}

func (e eventSaleOffRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventDiscount, error) {
	eventDiscountCollection := e.database.Collection(e.collectionEventDiscount)

	filter := bson.M{"_id": id}
	var eventDiscount domain.EventDiscount
	if err := eventDiscountCollection.FindOne(ctx, filter).Decode(&eventDiscount); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.EventDiscount{}, nil
		}
		return domain.EventDiscount{}, err
	}

	return eventDiscount, nil
}
func (e eventSaleOffRepository) GetByUserIDInApplicableAndExpiringOneDayLeft(ctx context.Context, userID primitive.ObjectID) (domain.EventDiscount, error) {
	eventDiscountCollection := e.database.Collection(e.collectionEventDiscount)

	filter := bson.M{
		"applicable_users": userID,
		"end_date":         bson.M{"$lt": time.Now().Add(24 * time.Hour)},
	}

	var eventDiscount domain.EventDiscount
	err := eventDiscountCollection.FindOne(ctx, filter).Decode(&eventDiscount)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.EventDiscount{}, nil
		}
		return domain.EventDiscount{}, err
	}

	return eventDiscount, nil
}

func (e eventSaleOffRepository) GetByUserIDInApplicableAndEventID(ctx context.Context, userID primitive.ObjectID, eventID primitive.ObjectID) (domain.EventDiscount, error) {
	eventDiscountCollection := e.database.Collection(e.collectionEventDiscount)

	filter := bson.M{"event_id": eventID, "applicable_users": userID}
	var eventDiscount domain.EventDiscount
	if err := eventDiscountCollection.FindOne(ctx, filter).Decode(&eventDiscount); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.EventDiscount{}, nil
		}
		return domain.EventDiscount{}, err
	}

	return eventDiscount, nil
}

func (e eventSaleOffRepository) GetByUserIDInApplicable(ctx context.Context, userID primitive.ObjectID) ([]domain.EventDiscount, error) {
	eventDiscountCollection := e.database.Collection(e.collectionEventDiscount)

	filter := bson.M{"applicable_users": userID}
	cursor, err := eventDiscountCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventDiscounts []domain.EventDiscount
	for cursor.Next(ctx) {
		var eventDiscount domain.EventDiscount
		if err = cursor.Decode(&eventDiscount); err != nil {
			return nil, err
		}

		eventDiscounts = append(eventDiscounts, eventDiscount)
	}

	return eventDiscounts, nil
}

func (e eventSaleOffRepository) GetAll(ctx context.Context) ([]domain.EventDiscount, error) {
	eventDiscountCollection := e.database.Collection(e.collectionEventDiscount)

	filter := bson.M{}
	cursor, err := eventDiscountCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventDiscounts []domain.EventDiscount
	for cursor.Next(ctx) {
		var eventDiscount domain.EventDiscount
		if err = cursor.Decode(&eventDiscount); err != nil {
			return nil, err
		}

		eventDiscounts = append(eventDiscounts, eventDiscount)
	}

	return eventDiscounts, nil
}

func (e eventSaleOffRepository) CreateOne(ctx context.Context, discount domain.EventDiscount) error {
	eventDiscountCollection := e.database.Collection(e.collectionEventDiscount)

	if err := validate_data.ValidateEventDiscount(discount); err != nil {
		return err
	}

	_, err := eventDiscountCollection.InsertOne(ctx, discount)
	if err != nil {
		return err
	}

	return nil
}

func (e eventSaleOffRepository) UpdateOne(ctx context.Context, discount domain.EventDiscount) error {
	eventDiscountCollection := e.database.Collection(e.collectionEventDiscount)

	if err := validate_data.ValidateEventDiscount(discount); err != nil {
		return err
	}

	filter := bson.M{"_id": discount.ID}
	update := bson.M{"$set": bson.M{
		"event_id":         discount.EventID,
		"discount_value":   discount.DiscountValue,
		"discount_unit":    discount.DiscountUnit,
		"date_created":     discount.DateCreated,
		"start_date":       discount.StartDate,
		"end_date":         discount.EndDate,
		"applicable_users": discount.ApplicableUsers,
	}}

	_, err := eventDiscountCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e eventSaleOffRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	eventDiscountCollection := e.database.Collection(e.collectionEventDiscount)

	filter := bson.M{"_id": id}
	_, err := eventDiscountCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
