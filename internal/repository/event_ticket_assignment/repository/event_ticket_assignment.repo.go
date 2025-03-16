package event_ticket_assignment_repository

import (
	"bookify/internal/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IEventTicketAssignmentRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventTicketAssignment, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.EventTicketAssignment, error)
	GetByEventID(ctx context.Context, eventID primitive.ObjectID) ([]domain.EventTicketAssignment, error)
	GetAll(ctx context.Context) ([]domain.EventTicketAssignment, error)
	CreateOne(ctx context.Context, eventTicketAssignment domain.EventTicketAssignment) error
	UpdateOne(ctx context.Context, eventTicketAssignment domain.EventTicketAssignment) error
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	StatisticsRevenueByEventID(ctx context.Context, eventId primitive.ObjectID) (domain.EventTicketAssignmentResponse, error)
}

type eventTicketAssignmentRepository struct {
	database                        *mongo.Database
	collectionEventTicketAssignment string
}

func NewEventTicketAssignmentRepository(database *mongo.Database, collectionEventTicketAssignment string) IEventTicketAssignmentRepository {
	return &eventTicketAssignmentRepository{database: database, collectionEventTicketAssignment: collectionEventTicketAssignment}
}

func (e *eventTicketAssignmentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.EventTicketAssignment, error) {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{"_id": id}
	var eventTicketAssignment domain.EventTicketAssignment
	if err := eventTicketAssignmentCollection.FindOne(ctx, filter).Decode(&eventTicketAssignment); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.EventTicketAssignment{}, nil
		}
		return domain.EventTicketAssignment{}, err
	}

	return eventTicketAssignment, nil
}

func (e *eventTicketAssignmentRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.EventTicketAssignment, error) {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{"user_id": userID}
	cursor, err := eventTicketAssignmentCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventTicketAssignments []domain.EventTicketAssignment
	for cursor.Next(ctx) {
		var eventTicketAssignment domain.EventTicketAssignment
		if err = cursor.Decode(&eventTicketAssignment); err != nil {
			return nil, err
		}

		eventTicketAssignments = append(eventTicketAssignments, eventTicketAssignment)
	}

	return eventTicketAssignments, nil
}

func (e *eventTicketAssignmentRepository) GetByEventID(ctx context.Context, eventID primitive.ObjectID) ([]domain.EventTicketAssignment, error) {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{"event_id": eventID}
	cursor, err := eventTicketAssignmentCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventTicketAssignments []domain.EventTicketAssignment
	for cursor.Next(ctx) {
		var eventTicketAssignment domain.EventTicketAssignment
		if err = cursor.Decode(&eventTicketAssignment); err != nil {
			return nil, err
		}

		eventTicketAssignments = append(eventTicketAssignments, eventTicketAssignment)
	}

	return eventTicketAssignments, nil
}

func (e *eventTicketAssignmentRepository) GetAll(ctx context.Context) ([]domain.EventTicketAssignment, error) {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{}
	cursor, err := eventTicketAssignmentCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventTicketAssignments []domain.EventTicketAssignment
	for cursor.Next(ctx) {
		var eventTicketAssignment domain.EventTicketAssignment
		if err = cursor.Decode(&eventTicketAssignment); err != nil {
			return nil, err
		}

		eventTicketAssignments = append(eventTicketAssignments, eventTicketAssignment)
	}

	return eventTicketAssignments, nil
}

func (e *eventTicketAssignmentRepository) CreateOne(ctx context.Context, eventTicketAssignment domain.EventTicketAssignment) error {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	// Need check duplicate
	_, err := eventTicketAssignmentCollection.InsertOne(ctx, eventTicketAssignment)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketAssignmentRepository) UpdateOne(ctx context.Context, eventTicketAssignment domain.EventTicketAssignment) error {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{"id": eventTicketAssignment.ID}
	update := bson.M{
		"$set": bson.M{
			"event_id":      eventTicketAssignment.EventID,
			"attendance_id": eventTicketAssignment.AttendanceID,
			"purchase_date": eventTicketAssignment.PurchaseDate,
			"expiry_date":   eventTicketAssignment.ExpiryDate,
			"price":         eventTicketAssignment.Price,
			"ticket_type":   eventTicketAssignment.TicketType,
			"status":        eventTicketAssignment.Status,
		},
	}

	_, err := eventTicketAssignmentCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketAssignmentRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	_, err := eventTicketAssignmentCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketAssignmentRepository) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	filter := bson.M{"id": id}
	_, err := eventTicketAssignmentCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (e *eventTicketAssignmentRepository) StatisticsRevenueByEventID(ctx context.Context, eventId primitive.ObjectID) (domain.EventTicketAssignmentResponse, error) {
	eventTicketAssignmentCollection := e.database.Collection(e.collectionEventTicketAssignment)

	// Sử dụng aggregation pipeline để nhóm và tính tổng revenue
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"event_id": eventId}}}, // Lọc theo event_id
		{
			{"$group", bson.M{
				"_id":          nil, // Không cần nhóm theo trường nào, chỉ tính tổng toàn bộ
				"totalRevenue": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$price", "$quantity"}}},
			}},
		},
	}

	cursor, err := eventTicketAssignmentCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return domain.EventTicketAssignmentResponse{}, err
	}
	defer cursor.Close(ctx)

	var totalRevenue float64
	if cursor.Next(ctx) {
		var result struct {
			TotalRevenue float64 `bson:"totalRevenue"`
		}
		if err := cursor.Decode(&result); err != nil {
			return domain.EventTicketAssignmentResponse{}, err
		}
		totalRevenue = result.TotalRevenue
	}

	// Kiểm tra lỗi con trỏ
	if err = cursor.Err(); err != nil {
		return domain.EventTicketAssignmentResponse{}, err
	}

	// Lấy danh sách vé chi tiết (nếu cần)
	var eventTicketAssignments []domain.EventTicketAssignment
	cursor2, err := eventTicketAssignmentCollection.Find(ctx, bson.M{"event_id": eventId})
	if err != nil {
		return domain.EventTicketAssignmentResponse{}, err
	}
	defer cursor2.Close(ctx)

	if err = cursor2.All(ctx, &eventTicketAssignments); err != nil {
		return domain.EventTicketAssignmentResponse{}, err
	}

	// Trả về response
	response := domain.EventTicketAssignmentResponse{
		EventTicketAssignment: eventTicketAssignments,
		TotalRevenue:          totalRevenue,
	}

	return response, nil
}
