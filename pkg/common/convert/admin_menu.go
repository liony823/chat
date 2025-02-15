package convert

import (
	"context"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	adminpb "github.com/openimsdk/chat/pkg/protocol/common"
)

func AdminMenuDB2PB(ctx context.Context, menus []*admindb.AdminMenu) (adminMenuPbs []*adminpb.AdminMenu, err error) {
	if len(menus) == 0 {
		return nil, nil
	}

	for _, menu := range menus {
		adminMenuPb := &adminpb.AdminMenu{}
		adminMenuPb.Key = menu.Key
		adminMenuPb.Path = menu.Path
		adminMenuPb.Name = menu.Name
		adminMenuPb.Icon = menu.Icon
		adminMenuPb.Sort = menu.Sort
		adminMenuPb.Parent = menu.Parent
		adminMenuPb.Layout = menu.Layout
		adminMenuPb.Component = menu.Component
		adminMenuPb.HiddenInMenu = menu.HiddenInMenu
		adminMenuPb.Redirect = menu.Redirect

		adminMenuPbs = append(adminMenuPbs, adminMenuPb)
	}
	return adminMenuPbs, nil
}
