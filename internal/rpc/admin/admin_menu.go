package admin

import (
	"context"
	"errors"
	"strings"

	"github.com/openimsdk/chat/pkg/common/convert"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"go.mongodb.org/mongo-driver/mongo"
)

func (o *adminServer) CreateAdminMenu(ctx context.Context, req *adminpb.CreateAdminMenuReq) (*adminpb.CreateAdminMenuResp, error) {
	create, err := ToDBAdminMenuCreate(req)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	err = o.Database.CreateAdminMenu(ctx, []*admindb.AdminMenu{create})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &adminpb.CreateAdminMenuResp{}, nil
}

func (o *adminServer) UpdateAdminMenu(ctx context.Context, req *adminpb.UpdateAdminMenuReq) (*adminpb.UpdateAdminMenuResp, error) {
	update, err := ToDBAdminMenuUpdate(req)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	err = o.Database.UpdateAdminMenu(ctx, req.Key, update)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &adminpb.UpdateAdminMenuResp{}, nil
}

func (o *adminServer) DeleteAdminMenu(ctx context.Context, req *adminpb.DeleteAdminMenuReq) (*adminpb.DeleteAdminMenuResp, error) {
	err := o.Database.DeleteAdminMenu(ctx, req.Keys)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &adminpb.DeleteAdminMenuResp{}, nil
}

func (o *adminServer) TakeAdminMenu(ctx context.Context, req *adminpb.TakeAdminMenuReq) (*adminpb.TakeAdminMenuResp, error) {
	menu, err := o.Database.TakeAdminMenu(ctx, req.Key)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	var resp adminpb.TakeAdminMenuResp
	if err := datautil.CopyStructFields(menu, &resp); err != nil {
		return nil, errs.Wrap(err)
	}
	return &resp, nil
}

func (o *adminServer) ListAdminMenu(ctx context.Context, req *adminpb.ListAdminMenuReq) (*adminpb.ListAdminMenuResp, error) {
	res, err := o.Database.ListAdminMenu(ctx, req.Parent)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	menus, err := convert.AdminMenuDB2PB(ctx, res)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &adminpb.ListAdminMenuResp{Menus: menus}, nil
}

// 获取当前用户的菜单
func (o *adminServer) ListAdminUserMenu(ctx context.Context, req *adminpb.ListAdminUserMenuReq) (*adminpb.ListAdminUserMenuResp, error) {
	var menus []*admindb.AdminMenu
	var err error

	// Super admin can access all menus
	if err := o.CheckSuperAdmin(ctx); err == nil {
		menus, err = o.Database.ListAdminMenu(ctx, "")
		if err != nil {
			return nil, errs.Wrap(err)
		}
	} else {
		userId, err := mctx.CheckAdmin(ctx)
		if err != nil {
			return nil, err
		}
		userMenu, err := o.Database.TakeAdminUserMenu(ctx, userId)
		if errors.Is(err, mongo.ErrNoDocuments) {
			userMenu = &admindb.AdminUserMenu{
				UserID: userId,
				Menus:  []string{},
			}
		} else if err != nil {
			return nil, errs.Wrap(err)
		}

		// Check if any menu key contains a dash but missing its parent
		menuMap := make(map[string]bool)
		for _, key := range userMenu.Menus {
			menuMap[key] = true
		}

		var newMenus []string
		for _, key := range userMenu.Menus {
			if idx := strings.Index(key, "-"); idx > 0 {
				parent := key[:idx]
				if !menuMap[parent] {
					newMenus = append(newMenus, parent)
					menuMap[parent] = true
				}
			}
		}

		userMenu.Menus = append(userMenu.Menus, newMenus...)

		menus, err = o.Database.ListAdminMenuByKeys(ctx, userMenu.Menus)
		if err != nil {
			return nil, errs.Wrap(err)
		}
	}

	// Convert to protobuf format
	menusPB, err := convert.AdminMenuDB2PB(ctx, menus)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return &adminpb.ListAdminUserMenuResp{Menus: menusPB}, nil
}

// 分配菜单给用户
func (o *adminServer) AssignAdminUserMenu(ctx context.Context, req *adminpb.AssignAdminUserMenuReq) (*adminpb.AssignAdminUserMenuResp, error) {
	_, err := o.Database.TakeAdminUserMenu(ctx, req.UserID)

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		record := &admindb.AdminUserMenu{
			UserID: req.UserID,
			Menus:  req.Menus,
		}

		err = o.Database.CreateAdminUserMenu(ctx, []*admindb.AdminUserMenu{record})
		if err != nil {
			return nil, errs.Wrap(err)
		}
		return &adminpb.AssignAdminUserMenuResp{}, nil
	case err != nil:
		return nil, errs.Wrap(err)
	}

	err = o.Database.UpdateAdminUserMenu(ctx, req.UserID, req.Menus)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &adminpb.AssignAdminUserMenuResp{}, nil
}

// 根据用户id获取菜单
func (o *adminServer) GetAdminUserMenu(ctx context.Context, req *adminpb.GetAdminUserMenuReq) (*adminpb.GetAdminUserMenuResp, error) {
	_, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	userMenu, err := o.Database.TakeAdminUserMenu(ctx, req.UserID)
	if errors.Is(err, mongo.ErrNoDocuments) {
		userMenu = &admindb.AdminUserMenu{
			UserID: req.UserID,
			Menus:  []string{},
		}
	} else if err != nil {
		return nil, errs.Wrap(err)
	}
	menus, err := o.Database.ListAdminMenuByKeys(ctx, userMenu.Menus)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	menusPB, err := convert.AdminMenuDB2PB(ctx, menus)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &adminpb.GetAdminUserMenuResp{Menus: menusPB}, nil
}
