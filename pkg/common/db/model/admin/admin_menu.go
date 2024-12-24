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

func NewAdminMenu(db *mongo.Database) (admin.AdminMenuInterface, error) {
	coll := db.Collection("admin_menu")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "key", Value: 1},
			{Key: "path", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &AdminMenu{coll: coll}, nil
}

type AdminMenu struct {
	coll *mongo.Collection
}

func (o *AdminMenu) Create(ctx context.Context, menus []*admin.AdminMenu) error {
	return mongoutil.InsertMany(ctx, o.coll, menus)
}

func (o *AdminMenu) Update(ctx context.Context, key string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"key": key}, bson.M{"$set": data}, false)
}

func (o *AdminMenu) Delete(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"key": bson.M{"$in": keys}})
}

func (o *AdminMenu) Take(ctx context.Context, key string) (*admin.AdminMenu, error) {
	return mongoutil.FindOne[*admin.AdminMenu](ctx, o.coll, bson.M{"key": key})
}

func (o *AdminMenu) TakeByParent(ctx context.Context, parent string) ([]*admin.AdminMenu, error) {
	return mongoutil.Find[*admin.AdminMenu](ctx, o.coll, bson.M{"parent": parent})
}

func (o *AdminMenu) List(ctx context.Context, parent string) ([]*admin.AdminMenu, error) {
	filter := bson.M{}
	if parent != "" {
		filter["parent"] = parent
	}
	return mongoutil.Find[*admin.AdminMenu](ctx, o.coll, filter)
}

func (o *AdminMenu) ListByKeys(ctx context.Context, keys []string) ([]*admin.AdminMenu, error) {
	return mongoutil.Find[*admin.AdminMenu](ctx, o.coll, bson.M{"key": bson.M{"$in": keys}})
}
