package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary 创建菜单
// @Description 创建菜单
// @Tags admin_menu
// @Id createAdminMenu
// @Accept json
// @Produce json
// @Param admin_menu body admin.CreateAdminMenuReq true "菜单信息"
// @Success 200 {object} apiresp.ApiResponse{data=admin.CreateAdminMenuResp}
// @Router /admin_menu/create [post]
func (o *Api) CreateAdminMenu(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.CreateAdminMenu, o.adminClient)
}

// @Summary 更新菜单
// @Description 更新菜单
// @Tags admin_menu
// @Id updateAdminMenu
// @Accept json
// @Produce json
// @Param admin_menu body admin.UpdateAdminMenuReq true "菜单信息"
// @Success 200 {object} apiresp.ApiResponse{data=admin.UpdateAdminMenuResp}
// @Router /admin_menu/update [post]
func (o *Api) UpdateAdminMenu(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.UpdateAdminMenu, o.adminClient)
}

// @Summary 删除菜单
// @Description 删除菜单
// @Tags admin_menu
// @Id deleteAdminMenu
// @Accept json
// @Produce json
// @Param admin_menu body admin.DeleteAdminMenuReq true "菜单信息"
// @Success 200 {object} apiresp.ApiResponse{data=admin.DeleteAdminMenuResp}
// @Router /admin_menu/delete [post]
func (o *Api) DeleteAdminMenu(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DeleteAdminMenu, o.adminClient)
}

// @Summary 获取菜单
// @Description 获取菜单
// @Tags admin_menu
// @Id takeAdminMenu
// @Accept json
// @Produce json
// @Param admin_menu body admin.TakeAdminMenuReq true "菜单信息"
// @Success 200 {object} apiresp.ApiResponse{data=admin.TakeAdminMenuResp}
// @Router /admin_menu/take [post]
func (o *Api) TakeAdminMenu(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.TakeAdminMenu, o.adminClient)
}

// @Summary 获取菜单列表
// @Description 获取菜单列表
// @Tags admin_menu
// @Id listAdminMenu
// @Accept json
// @Produce json
// @Param admin_menu body admin.ListAdminMenuReq true "菜单信息"
// @Success 200 {object} apiresp.ApiResponse{data=admin.ListAdminMenuResp}
// @Router /admin_menu/list [post]
func (o *Api) ListAdminMenu(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.ListAdminMenu, o.adminClient)
}

// @Summary 获取用户权限菜单
// @Description 获取用户权限菜单
// @Tags admin_menu
// @Id listAdminUserMenu
// @Accept json
// @Produce json
// @Success 200 {object} apiresp.ApiResponse{data=admin.ListAdminUserMenuResp}
// @Router /admin_menu/user_menu [post]
func (o *Api) ListAdminUserMenu(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.ListAdminUserMenu, o.adminClient)
}

// @Summary 获取用户权限菜单
// @Description 获取用户权限菜单
// @Tags admin_menu
// @Id getAdminUserMenu
// @Accept json
// @Produce json
// @Param admin_menu body admin.GetAdminUserMenuReq true "菜单信息"
// @Success 200 {object} apiresp.ApiResponse{data=admin.GetAdminUserMenuResp}
// @Router /admin_menu/user_menu/get [post]
func (o *Api) GetAdminUserMenu(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetAdminUserMenu, o.adminClient)
}

// @Summary 分配用户权限菜单
// @Description 分配用户权限菜单
// @Tags admin_menu
// @Id assignAdminUserMenu
// @Accept json
// @Produce json
// @Param admin_menu body admin.AssignAdminUserMenuReq true "菜单信息"
// @Success 200 {object} apiresp.ApiResponse{data=admin.AssignAdminUserMenuResp}
// @Router /admin_menu/user_menu/assign [post]
func (o *Api) AssignAdminUserMenu(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AssignAdminUserMenu, o.adminClient)
}
