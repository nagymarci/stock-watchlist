package database

import (
	"context"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Recommendations struct {
	collection *mongo.Collection
}

type recommendation struct {
	ID     primitive.ObjectID `bson:"_id"`
	Stocks []string           `bson:"stocks"`
}

func NewRecommendations(db *mongo.Database) *Recommendations {
	return &Recommendations{
		collection: db.Collection("recommendations"),
	}
}

func (r *Recommendations) Create(id primitive.ObjectID, stocks []string) error {
	_, err := r.collection.InsertOne(context.TODO(), recommendation{ID: id, Stocks: stocks})

	return err
}

func (r *Recommendations) Get(id primitive.ObjectID) ([]string, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	var result recommendation
	err := r.collection.FindOne(context.TODO(), filter).Decode(&result)

	return result.Stocks, err
}

func (r *Recommendations) Update(log *logrus.Logger, id primitive.ObjectID, stocks []string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	opts := options.Replace().SetUpsert(true)

	_, err := r.collection.ReplaceOne(context.TODO(), filter, recommendation{ID: id, Stocks: stocks}, opts)

	if err != nil {
		log.Errorln(err)
		return err
	}

	log.Infoln("recommendation inserted into DB")
	return nil
}
