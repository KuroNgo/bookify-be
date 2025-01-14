package repository

import (
	"bookify/internal/domain"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepository interface {
	FetchMany(ctx context.Context) ([]domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (domain.User, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (domain.User, error)

	UpdateOne(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, user *domain.User) error
	UpdateVerify(ctx context.Context, user *domain.User) error
	UpdateVerificationCode(ctx context.Context, user *domain.User) error
	UpsertOne(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateImage(ctx context.Context, user *domain.User) error

	UserExists(ctx context.Context, email string) (bool, error)
	CreateOne(ctx context.Context, user *domain.User) error
	DeleteOne(ctx context.Context, userID primitive.ObjectID) error
}

type userRepository struct {
	database       *mongo.Database
	collectionUser string
}

func (u userRepository) FetchMany(ctx context.Context) ([]domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{}
	cursor, err := collectionUser.Find(ctx, filter)
	if err != nil {
		return nil, errors.New(err.Error() + "error in the finding user into the database")
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var users []domain.User
	users = make([]domain.User, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		var user domain.User
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	// Check for any errors encountered during iteration
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"email": email}
	var user domain.User
	if err := collectionUser.FindOne(ctx, filter).Decode(&user); err != nil {
		return domain.User{}, errors.New(err.Error() + "error in the finding user into the database")
	}

	return user, nil
}

func (u userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"_id": id}
	projection := bson.M{"password_hash": 0}

	var user domain.User
	if err := collectionUser.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&user); err != nil {
		return domain.User{}, errors.New(err.Error() + "error in the finding user into the database")
	}

	return user, nil
}

func (u userRepository) GetByVerificationCode(ctx context.Context, verificationCode string) (domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"verification_code": verificationCode}
	var user domain.User
	if err := collectionUser.FindOne(ctx, filter).Decode(&user); err != nil {
		return domain.User{}, errors.New(err.Error() + "error in the finding user's data into database")
	}

	return user, nil
}

func (u userRepository) UserExists(ctx context.Context, email string) (bool, error) {
	collectionUser := u.database.Collection(u.collectionUser)
	filter := bson.M{"email": email}
	count, err := collectionUser.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u userRepository) UpdateOne(ctx context.Context, user *domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	exists, err := u.UserExists(ctx, user.Email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("user is confirmed to exist")
	}

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err = collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New(err.Error() + "error in the updating user into database")
	}

	return nil
}

func (u userRepository) UpdatePassword(ctx context.Context, user *domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{"password_hash": user.PasswordHash}}

	_, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New(err.Error() + "error in the updating user into database")
	}

	return nil
}

func (u userRepository) UpdateVerify(ctx context.Context, user *domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"verify":     user.Verified,
		"updated_at": user.UpdatedAt,
	}}

	_, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New(err.Error() + "error in the updating user's data into database")
	}

	return nil
}

func (u userRepository) UpdateVerificationCode(ctx context.Context, user *domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"verify":            user.Verified,
		"verification_code": user.VerificationCode,
		"updated_at":        user.UpdatedAt,
	}}

	_, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New(err.Error() + "error in the updating user's data into database")
	}

	return nil
}

func (u userRepository) UpsertOne(ctx context.Context, user *domain.User) (*domain.User, error) {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"email": user.Email}

	// Chuẩn bị các giá trị cập nhật
	update := bson.D{{Key: "$set", Value: bson.M{
		"full_name":  user.FullName,
		"email":      user.Email,
		"avatar_url": user.AvatarURL,
		"phone":      user.Phone,
		"provider":   user.Provider,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"role":       user.Role,
	}}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	result := collectionUser.FindOneAndUpdate(ctx, filter, update, opts)

	var updatedUser *domain.User
	if err := result.Decode(&updatedUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return updatedUser, nil
}

func (u userRepository) UpdateImage(ctx context.Context, user *domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"avatar_url": user.AvatarURL,
		"asset_url":  user.AssetURL,
	}}

	_, err := collectionUser.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New(err.Error() + "error in the updating user's data into database")
	}

	return nil
}

func (u userRepository) CreateOne(ctx context.Context, user *domain.User) error {
	collectionUser := u.database.Collection(u.collectionUser)

	exists, err := u.UserExists(ctx, user.Email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("user is confirmed to exist")
	}

	_, err = collectionUser.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error inserting user into the database: %w", err)
	}

	return nil
}

func (u userRepository) DeleteOne(ctx context.Context, userID primitive.ObjectID) error {
	collectionUser := u.database.Collection(u.collectionUser)

	filter := bson.M{"_id": userID}
	_, err := collectionUser.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New(err.Error() + "error the deleting user's data into the database")
	}

	return nil
}

func NewUserRepository(db *mongo.Database, collectionUser string) IUserRepository {
	return &userRepository{database: db, collectionUser: collectionUser}
}
