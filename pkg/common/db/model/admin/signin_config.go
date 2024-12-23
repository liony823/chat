package admin

import (
	"context"

	"github.com/liony823/tools/db/mongoutil"
	"github.com/liony823/tools/errs"
	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewSigninConfig(db *mongo.Database) (admin.SigninConfigInterface, error) {
	coll := db.Collection("signin_config")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "key", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &SigninConfig{
		coll: coll,
	}, nil
}

type SigninConfig struct {
	coll *mongo.Collection
}

func (o *SigninConfig) Set(ctx context.Context, config *admin.SigninConfig) error {
	_, err := mongoutil.UpdateMany(ctx, o.coll, bson.M{}, bson.M{"$set": config}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (o *SigninConfig) Get(ctx context.Context) (*admin.SigninConfig, error) {
	cs, err := mongoutil.Find[*admin.SigninConfig](ctx, o.coll, bson.M{})
	if err != nil {
		return nil, err
	}
	if len(cs) == 0 {
		return nil, errs.WrapMsg(errs.ErrRecordNotFound, "signin config not found")
	}
	return cs[0], nil
}
