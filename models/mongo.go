package models

import (
	"context"
	"fmt"
	"github.com/misgorod/tucktuck-pull/common"
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
