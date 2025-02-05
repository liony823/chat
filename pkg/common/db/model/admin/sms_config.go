package admin

import (
	"context"

	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewSmsConfig(db *mongo.Database) (admin.SmsConfigInterface, error) {
	coll := db.Collection("sms_config")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "key", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &SmsConfig{
		coll: coll,
	}, nil
}

type SmsConfig struct {
	coll *mongo.Collection
}

func (o *SmsConfig) Set(ctx context.Context, config map[string]interface{}) error {
	for key, value := range config {
		filter := bson.M{"key": key}
		update := bson.M{
			"options": value,
		}
		err := mongoutil.UpdateOne(ctx, o.coll, filter, bson.M{"$set": update}, false, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *SmsConfig) Get(ctx context.Context) (map[string]interface{}, error) {
	cs, err := mongoutil.Find[*admin.SmsConfig](ctx, o.coll, bson.M{})
	if err != nil {
		return nil, err
	}
	cm := make(map[string]interface{})
	for _, config := range cs {
		cm[config.Key] = config.Options
	}
	return cm, nil
}
