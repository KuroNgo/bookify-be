package activity_log_repository

import (
	"bookify/internal/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IActivityLogRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.ActivityLog, error)
	GetByLevel(ctx context.Context, level string) ([]domain.ActivityLog, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.ActivityLog, error)
	GetAll(ctx context.Context) ([]domain.ActivityLog, error)
	CreateOne(ctx context.Context, activityLog *domain.ActivityLog) error
	UpdateOne(ctx context.Context, activityLog *domain.ActivityLog) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
}

type activityLogRepository struct {
	database              *mongo.Database
	collectionActivityLog string
}

func NewActivityLogRepository(database *mongo.Database, collectionActivityLog string) IActivityLogRepository {
	return &activityLogRepository{database: database, collectionActivityLog: collectionActivityLog}
}

func (a *activityLogRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.ActivityLog, error) {
	activityCollection := a.database.Collection(a.collectionActivityLog)

	filter := bson.M{"_id": id, "status": "enabled"}
	var log domain.ActivityLog
	if err := activityCollection.FindOne(ctx, filter).Decode(&log); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ActivityLog{}, nil
		}
		return domain.ActivityLog{}, err
	}

	return log, nil
}

func (a *activityLogRepository) GetByLevel(ctx context.Context, level string) ([]domain.ActivityLog, error) {
	activityCollection := a.database.Collection(a.collectionActivityLog)

	filter := bson.M{"level": level}
	cursor, err := activityCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []domain.ActivityLog
	for cursor.Next(ctx) {
		var log domain.ActivityLog
		if err = cursor.Decode(&log); err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (a *activityLogRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.ActivityLog, error) {
	activityCollection := a.database.Collection(a.collectionActivityLog)

	filter := bson.M{"user_id": userID}
	cursor, err := activityCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []domain.ActivityLog
	for cursor.Next(ctx) {
		var log domain.ActivityLog
		if err = cursor.Decode(&log); err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (a *activityLogRepository) GetAll(ctx context.Context) ([]domain.ActivityLog, error) {
	activityCollection := a.database.Collection(a.collectionActivityLog)

	filter := bson.M{}
	cursor, err := activityCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []domain.ActivityLog
	for cursor.Next(ctx) {
		var log domain.ActivityLog
		if err = cursor.Decode(&log); err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (a *activityLogRepository) CreateOne(ctx context.Context, activityLog *domain.ActivityLog) error {
	activityCollection := a.database.Collection(a.collectionActivityLog)

	_, err := activityCollection.InsertOne(ctx, activityLog)
	if err != nil {
		return err
	}

	return nil
}

func (a *activityLogRepository) UpdateOne(ctx context.Context, activityLog *domain.ActivityLog) error {
	activityCollection := a.database.Collection(a.collectionActivityLog)

	filter := bson.M{"_id": activityLog.ID}
	update := bson.M{"$set": bson.M{
		"client_ip":     activityLog.ClientIP,
		"user_id":       activityLog.UserID,
		"level":         activityLog.Level,
		"method":        activityLog.Method,
		"status_code":   activityLog.StatusCode,
		"body_size":     activityLog.BodySize,
		"path":          activityLog.Path,
		"latency":       activityLog.Latency,
		"activity_time": activityLog.ActivityTime,
		"expire_at":     activityLog.ExpireAt,
	}}

	_, err := activityCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *activityLogRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	activityCollection := a.database.Collection(a.collectionActivityLog)

	filter := bson.M{"_id": id}
	_, err := activityCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
