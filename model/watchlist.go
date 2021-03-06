package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Watchlist struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Name   string             `bson:"name" json:"name"`
	Stocks []string           `bson:"stocks" json:"stocks"`
	UserID string             `bson:"userId" json:"userId"`
}

type WatchlistRequest struct {
	Name   string   `bson:"name" json:"name"`
	Stocks []string `bson:"stocks" json:"stocks"`
	UserID string   `bson:"userId"`
}
