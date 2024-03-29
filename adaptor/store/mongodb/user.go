package store

import (
	"context"
	"time"

	"github.com/mohammadmokh/writino/adaptor/store/mongodb/models"
	"github.com/mohammadmokh/writino/entity"
	"github.com/mohammadmokh/writino/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m MongodbStore) CreateUser(ctx context.Context, user entity.User) error {

	coll := m.db.Collection("users")

	dbModel := models.MapFromUserEntity(user)
	dbModel.Id = primitive.NewObjectID()
	dbModel.CreatedAt = time.Now()
	_, err := coll.InsertOne(ctx, dbModel)

	if mongo.IsDuplicateKeyError(err) {
		return errors.ErrDuplicateRecord

	}
	return err
}

func (m MongodbStore) FindUser(ctx context.Context, id string) (entity.User, error) {

	coll := m.db.Collection("users")

	var user models.User

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.User{}, err
	}
	filter := bson.D{{"_id", objID}}
	res := coll.FindOne(ctx, filter, nil)

	if res.Err() == mongo.ErrNoDocuments {
		return entity.User{}, errors.ErrNotFound
	}

	err = res.Decode(&user)
	return models.MapToUserEntity(user), err
}

func (m MongodbStore) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {

	coll := m.db.Collection("users")

	var user models.User

	filter := bson.D{{"email", email}}
	res := coll.FindOne(ctx, filter, nil)

	if res.Err() == mongo.ErrNoDocuments {
		return entity.User{}, errors.ErrNotFound
	}

	err := res.Decode(&user)
	return models.MapToUserEntity(user), err
}

func (m MongodbStore) UpdateUser(ctx context.Context, user entity.User) error {

	coll := m.db.Collection("users")
	dbModel := models.MapFromUserEntity(user)

	filter := bson.D{{"_id", dbModel.Id}}
	_, err := coll.ReplaceOne(ctx, filter, dbModel, nil)

	if mongo.IsDuplicateKeyError(err) {
		return errors.ErrDuplicateRecord
	}

	return err
}

func (m MongodbStore) DeleteUser(ctx context.Context, id string) error {

	coll := m.db.Collection("users")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objID}}
	_, err = coll.DeleteOne(ctx, filter, nil)
	return err
}
