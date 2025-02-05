package admin

import (
	"context"

	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAdminUserMenu(db *mongo.Database) (admin.AdminUserMenuInterface, error) {
	coll := db.Collection("admin_user_menu")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &AdminUserMenu{coll: coll}, nil
}

type AdminUserMenu struct {
	coll *mongo.Collection
}

func (o *AdminUserMenu) Create(ctx context.Context, userMenus []*admin.AdminUserMenu) error {
	return mongoutil.InsertMany(ctx, o.coll, userMenus)
}

func (o *AdminUserMenu) Update(ctx context.Context, userID string, menus []string) error {
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"menus": menus}}, false)
}

func (o *AdminUserMenu) Delete(ctx context.Context, userIDs []string) error {
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (o *AdminUserMenu) Take(ctx context.Context, userID string) (*admin.AdminUserMenu, error) {
	return mongoutil.FindOne[*admin.AdminUserMenu](ctx, o.coll, bson.M{"user_id": userID})
}

func (o *AdminUserMenu) List(ctx context.Context) ([]*admin.AdminUserMenu, error) {
	return mongoutil.Find[*admin.AdminUserMenu](ctx, o.coll, bson.M{})
}
