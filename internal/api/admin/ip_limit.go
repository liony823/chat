package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary		添加用户IP登录限制
// @Description	限制用户在指定IP地址登录
// @Tags			forbidden
// @Id	addUserIPLimitLogin
// @Accept			json
// @Produce		json
// @Param			data	body		admin.AddUserIPLimitLoginReq	true	"IP限制信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/forbidden/user/add [post]
func (o *Api) AddUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddUserIPLimitLogin, o.adminClient, c)
}

// @Summary		搜索用户IP登录限制
// @Description	查询用户IP登录限制列表
// @Tags			ip_limit
// @Id	searchUserIPLimitLogin
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SearchUserIPLimitLoginReq	true	"搜索条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.SearchUserIPLimitLoginResp}
// @Router			/forbidden/user/search [post]
func (o *Api) SearchUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchUserIPLimitLogin, o.adminClient, c)
}

// @Summary		删除用户IP登录限制
// @Description	删除指定的用户IP登录限制
// @Tags			ip_limit
// @Id	delUserIPLimitLogin
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DelUserIPLimitLoginReq	true	"要删除的限制信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/forbidden/user/del [post]
func (o *Api) DelUserIPLimitLogin(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelUserIPLimitLogin, o.adminClient, c)
}

// @Summary		添加IP黑名单
// @Description	添加禁止注册和登录的IP
// @Tags			ip_limit
// @Id	addIPForbidden
// @Accept			json
// @Produce		json
// @Param			data	body		admin.AddIPForbiddenReq	true	"IP黑名单信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/forbidden/ip/add [post]
func (o *Api) AddIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddIPForbidden, o.adminClient, c)
}

// @Summary		删除IP黑名单
// @Description	删除禁止注册和登录的IP
// @Tags			ip_limit
// @Id	delIPForbidden
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DelIPForbiddenReq	true	"要删除的IP信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/forbidden/ip/del [post]
func (o *Api) DelIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelIPForbidden, o.adminClient, c)
}

// @Summary		搜索IP黑名单
// @Description	查询被禁止注册和登录的IP列表
// @Tags			ip_limit
// @Id	searchIPForbidden
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SearchIPForbiddenReq	true	"搜索条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.SearchIPForbiddenResp}
// @Router			/forbidden/ip/search [post]
func (o *Api) SearchIPForbidden(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchIPForbidden, o.adminClient, c)
}
