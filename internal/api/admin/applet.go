package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary		添加小程序
// @Description	添加新的小程序配置
// @Tags			applet
// @Id	addApplet
// @Accept			json
// @Produce		json
// @Param			data	body		admin.AddAppletReq	true	"小程序信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/applet/add [post]
func (o *Api) AddApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddApplet, o.adminClient, c)
}

// @Summary		删除小程序
// @Description	删除指定的小程序
// @Tags			applet
// @Id	delApplet
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DelAppletReq	true	"要删除的小程序ID"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/applet/del [post]
func (o *Api) DelApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelApplet, o.adminClient, c)
}

// @Summary		更新小程序
// @Description	更新小程序配置信息
// @Tags			applet
// @Id	updateApplet
// @Accept			json
// @Produce		json
// @Param			data	body		admin.UpdateAppletReq	true	"更新的小程序信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/applet/update [post]
func (o *Api) UpdateApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.UpdateApplet, o.adminClient, c)
}

// @Summary		搜索小程序
// @Description	根据条件搜索小程序列表
// @Tags			applet
// @Id	searchApplet
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SearchAppletReq	true	"搜索条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.SearchAppletResp}
// @Router			/applet/search [post]
func (o *Api) SearchApplet(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchApplet, o.adminClient, c)
}