package admin

import (
	"context"
	"time"

	"github.com/google/uuid"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) CreateAnnouncement(ctx context.Context, req *admin.CreateAnnouncementReq) (*admin.CreateAnnouncementResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	an := admindb.Announcement{
		AnnouncementID: uuid.New().String(),
		Title:          req.Title,
		Content:        req.Content,
		AppVersion:     req.AppVersion,
		AppLang:        req.AppLang,
		CreatedAt:      time.Now(),
	}
	if err := o.Database.CreateAnnouncement(ctx, []*admindb.Announcement{&an}); err != nil {
		return nil, err
	}
	return &admin.CreateAnnouncementResp{}, nil
}

func (o *adminServer) UpdateAnnouncement(ctx context.Context, req *admin.UpdateAnnouncementReq) (*admin.UpdateAnnouncementResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}

	update := make(map[string]any)

	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Content != "" {
		update["content"] = req.Content
	}
	if req.AppVersion != "" {
		update["app_version"] = req.AppVersion
	}
	if req.AppLang != "" {
		update["app_lang"] = req.AppLang
	}

	update["updated_at"] = time.Now()

	if err := o.Database.UpdateAnnouncement(ctx, req.AnnouncementID, update); err != nil {
		return nil, err
	}
	return &admin.UpdateAnnouncementResp{}, nil
}

func (o *adminServer) DeleteAnnouncement(ctx context.Context, req *admin.DeleteAnnouncementReq) (*admin.DeleteAnnouncementResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	count, err := o.Database.DelAnnouncement(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &admin.DeleteAnnouncementResp{
		DeletedCount: count,
	}, nil
}

func (o *adminServer) SearchAnnouncement(ctx context.Context, req *admin.SearchAnnouncementReq) (*admin.SearchAnnouncementResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, list, err := o.Database.SearchAnnouncement(ctx, req.Keyword, req.Status, req.AppVersion, req.Pagination)
	if err != nil {
		return nil, err
	}

	announcements := make([]*admin.Announcement, 0, len(list))
	for _, announcement := range list {
		announcements = append(announcements, &admin.Announcement{
			AnnouncementID: announcement.AnnouncementID,
			Title:          announcement.Title,
			Content:        announcement.Content,
			AppVersion:     announcement.AppVersion,
			AppLang:        announcement.AppLang,
			CreatedAt:      announcement.CreatedAt.Unix(),
			UpdatedAt:      announcement.UpdatedAt.Unix(),
			PublishedAt:    announcement.PublishedAt.Unix(),
			PublisherID:    announcement.PublisherID,
			Status:         announcement.Status,
		})
	}

	return &admin.SearchAnnouncementResp{
		Total:         total,
		Announcements: announcements,
	}, nil
}

func (o *adminServer) PublishAnnouncement(ctx context.Context, req *admin.PublishAnnouncementReq) (*admin.PublishAnnouncementResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}

	if err := o.Database.PublishAnnouncement(ctx, req.AnnouncementID); err != nil {
		return nil, err
	}
	return &admin.PublishAnnouncementResp{}, nil
}

func (o *adminServer) LatestAnnouncement(ctx context.Context, req *admin.LatestAnnouncementReq) (*admin.LatestAnnouncementResp, error) {
	if _, err := mctx.CheckUser(ctx); err != nil {
		return nil, err
	}
	announcement, err := o.Database.LatestAnnouncement(ctx, req.AppVersion, req.AppLang)
	if err != nil {
		return nil, err
	}

	return &admin.LatestAnnouncementResp{Announcement: &admin.Announcement{
		AnnouncementID: announcement.AnnouncementID,
		Title:          announcement.Title,
		Content:        announcement.Content,
		AppVersion:     announcement.AppVersion,
		AppLang:        announcement.AppLang,
		PublisherID:    announcement.PublisherID,
		CreatedAt:      announcement.CreatedAt.Unix(),
		UpdatedAt:      announcement.UpdatedAt.Unix(),
		PublishedAt:    announcement.PublishedAt.Unix(),
	}}, nil
}
