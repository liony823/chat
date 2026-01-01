package admin

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/pagination"
)

type Announcement struct {
	AnnouncementID string    `bson:"announcement_id"`
	PublisherID    string    `bson:"publisher_id"`
	IsDeleted      int32     `bson:"is_deleted"` //0: falase 1: true
	Status         int32     `bson:"status"`     //0: 草稿 1: 发布 2: 已过期
	Title          string    `bson:"title"`
	AppVersion     string    `bson:"app_version"`
	AppLang        string    `bson:"app_lang"`
	Content        string    `bson:"content"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
	PublishedAt    time.Time `bson:"published_at"`
}

func (Announcement) TableName() string {
	return "announcement"
}

type AnnouncementInterface interface {
	Update(ctx context.Context, id string, data map[string]any) error
	Create(ctx context.Context, ans []*Announcement) error
	Del(ctx context.Context, ids []string) (int64, error)
	Take(ctx context.Context, id string) (*Announcement, error)
	Find(ctx context.Context, app_version string, appLang string) ([]*Announcement, error)
	Search(ctx context.Context, keyword string, status int32, app_version string, pagination pagination.Pagination) (int64, []*Announcement, error)
	Latest(ctx context.Context, app_version string, app_lang string) (*Announcement, error)
	UpdateMany(ctx context.Context, ids []string, app_version string, app_lang string, data map[string]any) error
}
