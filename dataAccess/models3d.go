package dataaccess

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"ships3d/models"
)

// ErrTooManyModels is returned when a user already has MaxModelsPerUser models.
var ErrTooManyModels = errors.New("too many models")

// InsertModel stores an imported model. If the user already imported the same
// file (same hash), that one is returned instead and nothing is stored, so
// importing a file again doesn't count towards the limit.
func (s *Store) InsertModel(model *models.Model3d) (*models.Model3d, error) {
	ctx, cancel := newContext()
	defer cancel()
	collection := s.collection(modelsCollection)
	var existing models.Model3d
	err := collection.FindOne(ctx,
		bson.D{{Key: "userId", Value: model.UserId}, {Key: "hash", Value: model.Hash}},
		options.FindOne().SetProjection(bson.D{{Key: "data", Value: 0}}),
	).Decode(&existing)
	if err == nil {
		return &existing, nil
	}
	if !errors.Is(notFound(err), ErrNotFound) {
		return nil, err
	}
	count, err := collection.CountDocuments(ctx, bson.D{{Key: "userId", Value: model.UserId}})
	if err != nil {
		return nil, err
	}
	if count >= models.MaxModelsPerUser {
		return nil, ErrTooManyModels
	}
	model.Id = bson.NilObjectID
	res, err := collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}
	model.Id, _ = res.InsertedID.(bson.ObjectID)
	return model, nil
}

// GetModelData returns a model's file. Models are public: anyone who sees a
// ship needs its models.
func (s *Store) GetModelData(id bson.ObjectID) ([]byte, error) {
	ctx, cancel := newContext()
	defer cancel()
	var model models.Model3d
	err := s.collection(modelsCollection).FindOne(ctx, bson.D{{Key: "_id", Value: id}},
		options.FindOne().SetProjection(bson.D{{Key: "data", Value: 1}}),
	).Decode(&model)
	if err != nil {
		return nil, notFound(err)
	}
	return model.Data, nil
}
