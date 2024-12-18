package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/liony823/tools/apiresp"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary		封禁用户
// @Description	封禁指定用户账号
// @Tags			block
// @Id	addBlockUser
// @Accept			json
// @Produce		json
// @Param			data	body		admin.BlockUserReq	true	"要封禁的用户信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/block/add [post]
func (o *Api) BlockUser(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.BlockUserReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := o.adminClient.BlockUser(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	err = o.imApiCaller.ForceOffLine(mctx.WithApiToken(c, imToken), req.UserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// @Summary		解封用户
// @Description	解除用户账号封禁
// @Tags			block
// @Id	unblockUser
// @Accept			json
// @Produce		json
// @Param			data	body		admin.UnblockUserReq	true	"要解封的用户信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/block/del [post]
func (o *Api) UnblockUser(c *gin.Context) {
	a2r.Call(admin.AdminClient.UnblockUser, o.adminClient, c)
}

// @Summary		搜索被封禁用户
// @Description	查询被封禁的用户列表
// @Tags			block
// @Id	searchBlockUser
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SearchBlockUserReq	true	"搜索条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.SearchBlockUserResp}
// @Router			/block/search [post]
func (o *Api) SearchBlockUser(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchBlockUser, o.adminClient, c)
}
