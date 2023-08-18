package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"broabroad/internal/app/config"
)

const (
	seekRequestsCollection = "seekRequests"
)

type Repository interface {
	CreateSeekRequest(ctx context.Context, sr SeekRequest) error
	UpdateSeekRequestByID(ctx context.Context, id int) error
}

type Repo struct {
	mongoClient            *mongo.Client
	seekRequestsCollection *mongo.Collection
}

func NewRepo(ctx context.Context, dbCfg config.Database) (*Repo, error) {
	mongoClientOpts := options.Client().ApplyURI(dbCfg.ConnectionString)
	mongoClient, err := mongo.Connect(ctx, mongoClientOpts)
	if err != nil {
		return nil, fmt.Errorf("connect to Mongo DB: %w", err)
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("check connection to Mongo DB: %w", err)
	}
	return &Repo{
		mongoClient:            mongoClient,
		seekRequestsCollection: mongoClient.Database(dbCfg.Name).Collection(seekRequestsCollection),
	}, nil
}

type SeekRequest struct {
	ID           int    `json:"id"`
	DemandRating int    `json:"demand_rating"`
	RecipientID  int    `json:"recipient_id"`
	Street       string `json:"street"`
	PostCode     int    `json:"post_code"`
	City         string `json:"city"`
	Country      string `json:"country"`
}

func (r *Repo) CreateSeekRequest(ctx context.Context, sr SeekRequest) error {
	_, err := r.seekRequestsCollection.InsertOne(ctx, sr)
	if err != nil {
		return fmt.Errorf("create document: %w", err)
	}
	return nil
}

func (r *Repo) UpdateSeekRequestByID(ctx context.Context, id int) error {
	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$inc", bson.D{
			{"demand_rating", 1},
		}},
	}
	res, err := r.seekRequestsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update document: %w", err)
	}
	log.Printf("Update SeekRequest document: matched %v, updated %v", res.MatchedCount, res.ModifiedCount)
	return nil
}
