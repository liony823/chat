package admin

import (
	"context"
	"errors"
	"strings"

	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAnnouncement(db *mongo.Database) (admin.AnnouncementInterface, error) {
	coll := db.Collection("announcement")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "announcement_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &AnnouncementMgo{
		coll: coll,
	}, nil
}

type AnnouncementMgo struct {
	coll *mongo.Collection
}

// Create implements admin.AnnouncementInterface.
func (a *AnnouncementMgo) Create(ctx context.Context, ans []*admin.Announcement) error {
	return mongoutil.InsertMany(ctx, a.coll, ans)
}

// Del implements admin.AnnouncementInterface.
func (a *AnnouncementMgo) Del(ctx context.Context, ids []string) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	res, err := mongoutil.UpdateMany(ctx, a.coll, bson.M{"announcement_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"is_deleted": 1}})
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

// Latest implements admin.AnnouncementInterface.
func (a *AnnouncementMgo) Latest(ctx context.Context, app_version string, app_lang string) (*admin.Announcement, error) {
	if len(app_version) == 0 || len(app_lang) == 0 {
		return nil, errors.New("invalid version or lang")
	}

	return mongoutil.FindOne[*admin.Announcement](ctx, a.coll, bson.M{"app_version": app_version, "app_lang": strings.ToLower(app_lang), "status": 1, "is_deleted": 0})
}

// Search implements admin.AnnouncementInterface.
func (a *AnnouncementMgo) Search(ctx context.Context, keyword string, status int32, app_version string, pagination pagination.Pagination) (int64, []*admin.Announcement, error) {
	filter := bson.M{
		"is_deleted": 0,
	}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": keyword, "$options": "i"}},
			{"content": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	if status != -1 {
		filter["status"] = status
	}

	if app_version != "" {
		filter["app_version"] = app_version
	}

	return mongoutil.FindPage[*admin.Announcement](ctx, a.coll, filter, pagination, options.Find().SetSort(bson.D{{Key: "published_at", Value: -1}, {Key: "updated_at", Value: -1}, {Key: "created_at", Value: -1}}))
}

// Take implements admin.AnnouncementInterface.
func (a *AnnouncementMgo) Take(ctx context.Context, id string) (*admin.Announcement, error) {
	return mongoutil.FindOne[*admin.Announcement](ctx, a.coll, bson.M{"announcement_id": id})
}

// Update implements admin.AnnouncementInterface.
func (a *AnnouncementMgo) Update(ctx context.Context, id string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, a.coll, bson.M{"announcement_id": id}, bson.M{"$set": data}, false)
}

func (a *AnnouncementMgo) Find(ctx context.Context, app_version string, app_lang string) ([]*admin.Announcement, error) {
	return mongoutil.Find[*admin.Announcement](ctx, a.coll, bson.M{"app_version": app_version, "app_lang": app_lang, "is_deleted": 0})
}

func (a *AnnouncementMgo) UpdateMany(ctx context.Context, ids []string, app_version string, app_lang string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	_, err := mongoutil.UpdateMany(ctx, a.coll, bson.M{"announcement_id": bson.M{"$in": ids}, "app_version": app_version, "app_lang": app_lang}, bson.M{"$set": data})
	if err != nil {
		return err
	}
	return nil
}
