package repository

import (
	"context"
	"fmt"
	"github.com/misgorod/tucktuck-pull/common"
	"github.com/misgorod/tucktuck-pull/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Client struct {
	*mongo.Collection
}

func New() (*Client, error) {
	host := common.GetEnv("MONGO_HOST", "mongo")
	port := common.GetEnv("MONGO_PORT", "27017")
	database := common.GetEnv("MONGO_DB", "tucktuck")
	collection := common.GetEnv("MONGO_COLL", "events")
	client, err := mongo.NewClient(options.Client().
		SetAppName("pull").
		SetConnectTimeout(time.Minute).
		SetHosts([]string{fmt.Sprintf("%s:%s", host, port)}),
	)
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}
	return &Client{
		client.Database(database).Collection(collection),
	}, nil
}

func (c *Client) UpsertMany(ctx context.Context, results []models.Result) (models.UpsertResult, error) {
	updateModels := make([]mongo.WriteModel, 0, len(results))
	for _, result := range results {
		model := mongo.NewUpdateOneModel().
			SetFilter(
				bson.D{{
					"_id", result.Id,
				}},
			).
			SetUpdate(
				bson.D{
					{"$set", result},
				},
			).
			SetUpsert(true)
		updateModels = append(updateModels, model)
	}
	bulkResult, err := c.BulkWrite(ctx, updateModels)
	if err != nil {
		return models.UpsertResult{}, err
	}
	return models.UpsertResult{
		MatchedCount:  bulkResult.MatchedCount,
		ModifiedCount: bulkResult.ModifiedCount,
		UpsertedCount: bulkResult.UpsertedCount,
		UpsertedID:    bulkResult.UpsertedIDs,
		LastTime:      time.Now(),
	}, nil
}
