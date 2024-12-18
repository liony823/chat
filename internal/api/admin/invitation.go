package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary		添加邀请码
// @Description	手动添加新的邀请码
// @Tags			invitation
// @Id	addInvitationCode
// @Accept			json
// @Produce		json
// @Param			data	body		admin.AddInvitationCodeReq	true	"邀请码信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/invitation_code/add [post]
func (o *Api) AddInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.AddInvitationCode, o.adminClient, c)
}

// @Summary		生成邀请码
// @Description	批量生成新的邀请码
// @Tags			invitation
// @Id	genInvitationCode
// @Accept			json
// @Produce		json
// @Param			data	body		admin.GenInvitationCodeReq	true	"生成邀请码的配置信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/invitation_code/gen [post]
func (o *Api) GenInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.GenInvitationCode, o.adminClient, c)
}

// @Summary		删除邀请码
// @Description	删除指定的邀请码
// @Tags			invitation
// @Id	delInvitationCode
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DelInvitationCodeReq	true	"要删除的邀请码"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/invitation_code/del [post]
func (o *Api) DelInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelInvitationCode, o.adminClient, c)
}

// @Summary		搜索邀请码
// @Description	根据条件搜索邀请码
// @Tags			invitation
// @Id	searchInvitationCode
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SearchInvitationCodeReq	true	"搜索条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.SearchInvitationCodeResp}
// @Router			/invitation_code/search [post]
func (o *Api) SearchInvitationCode(c *gin.Context) {
	a2r.Call(admin.AdminClient.SearchInvitationCode, o.adminClient, c)
}
