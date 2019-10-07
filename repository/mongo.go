package repository

import (
	"context"
	"fmt"
	"github.com/misgorod/tucktuck-pull/common"
	"github.com/misgorod/tucktuck-pull/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Client struct {
	*mongo.Collection
}

func New() (*Client, error) {
	host := common.GetEnv("MONGO_HOST", "localhost")
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
	return &Client{
		client.Database(database).Collection(collection),
	}, nil
}

func (c *Client) UpsertMany(ctx context.Context, results []models.Result) (models.UpsertResult, error) {
	updateModels := make([]mongo.WriteModel, len(results))
	for _, result := range results {
		update, err := bson.Marshal(result)
		if err != nil {
			return models.UpsertResult{}, err
		}
		model := mongo.NewUpdateOneModel().
			SetFilter(
				bson.M{
					"_id": result.Id,
				},
			).
			SetUpdate(
				bson.M{
					"$set": update,
				},
			).
			SetUpsert(true)
		updateModels = append(updateModels, model)
	}
	bulkResult, err := c.BulkWrite(ctx, updateModels)
	for _, model := range updateModels {
		log.WithField("updateModels", fmt.Sprintf("%+v", model)).Info()
	}
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
