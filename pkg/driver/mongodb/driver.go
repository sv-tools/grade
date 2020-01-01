package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/sv-go-tools/grade/pkg/driver"
	"github.com/sv-go-tools/grade/pkg/driver/json"
)

func Execute(cfg *driver.Config) error {
	if cfg.ConnectionURL == "" {
		return json.Execute(cfg)
	}

	client, err := buildClient(cfg)
	if err != nil {
		return err
	}

	collection := client.Database(cfg.Database).Collection(cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	_, err = collection.InsertOne(ctx, cfg)
	if err != nil {
		return err
	}
	return nil
}

func buildClient(cfg *driver.Config) (*mongo.Client, error) {
	if cfg.Database == "" {
		cs, err := connstring.Parse(cfg.ConnectionURL)
		if err != nil {
			return nil, err
		}
		cfg.Database = cs.Database
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.ConnectionURL))
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}
