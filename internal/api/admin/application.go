package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/a2r"
)

// @Summary		获取最新应用版本
// @Description	获取应用程序的最新版本信息
// @Tags			application
// @Id	latestVersion
// @Accept			json
// @Produce		json
// @Param			data	body		admin.LatestApplicationVersionReq	true	"查询条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.LatestApplicationVersionResp}
// @Router			/application/latest_version [post]
func (o *Api) LatestApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.LatestApplicationVersion, o.adminClient)
}

// @Summary		分页查询应用版本
// @Description	分页获取应用程序版本列表
// @Tags			application
// @Id	pageVersions
// @Accept			json
// @Produce		json
// @Param			data	body		admin.PageApplicationVersionReq	true	"分页查询条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.PageApplicationVersionResp}
// @Router			/application/page_versions [post]
func (o *Api) PageApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.PageApplicationVersion, o.adminClient)
}

// @Summary		添加应用版本
// @Description	添加新的应用程序版本信息
// @Tags			application
// @Id	addVersion
// @Accept			json
// @Produce		json
// @Param			data	body		admin.AddApplicationVersionReq	true	"版本信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/application/add_version [post]
func (o *Api) AddApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddApplicationVersion, o.adminClient)
}

// @Summary		更新应用版本
// @Description	更新应用程序版本信息
// @Tags			application
// @Id	updateVersion
// @Accept			json
// @Produce		json
// @Param			data	body		admin.UpdateApplicationVersionReq	true	"更新的版本信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/application/update_version [post]
func (o *Api) UpdateApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.UpdateApplicationVersion, o.adminClient)
}

// @Summary		删除应用版本
// @Description	删除指定的应用程序版本
// @Tags			application
// @Id	deleteVersion
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DeleteApplicationVersionReq	true	"要删除的版本信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/application/delete_version [post]
func (o *Api) DeleteApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DeleteApplicationVersion, o.adminClient)
}
