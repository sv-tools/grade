package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/sv-go-tools/grade/pkg/driver"
)

func Execute(cfg *driver.Config) error {
	if cfg.ConnectionURL == "" {
		data, err := json.Marshal(makeDocuments(cfg.Records))
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	client, err := buildClient(cfg)
	if err != nil {
		return err
	}

	collection := client.Database(cfg.Database).Collection(cfg.Collection)

	res, err := collection.InsertMany(context.Background(), makeDocuments(cfg.Records))
	if err != nil {
		return err
	}

	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: bson.M{"$in": res.InsertedIDs},
		},
		bson.E{
			Key:   "timestamp",
			Value: bson.M{"$exists": false},
		},
	}
	update := bson.D{
		bson.E{
			Key: "$currentDate",
			Value: bson.D{
				bson.E{
					Key:   "timestamp",
					Value: true,
				},
			},
		},
	}
	_, err = collection.UpdateMany(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return err
	}
	for _, objID := range res.InsertedIDs {
		fmt.Println(objID)
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

func makeDocuments(records driver.Records) []interface{} {
	var documents []interface{}
	for _, rec := range records {
		delete(rec, "measured")
		documents = append(documents, rec)
	}
	return documents
}
