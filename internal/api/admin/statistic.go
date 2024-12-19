package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/liony823/tools/apiresp"
	"github.com/openimsdk/chat/pkg/common/apistruct"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/protocol/user"
)

// @Summary 获取系统统计数据
// @Description 获取系统统计数据
// @Tags			statistic
// @Id	LoginRecord
// @Accept			json
// @Produce		json
// @Success		200	{object}	apiresp.ApiResponse{data=admin.GetUserLoginRecordResp}
// @Router			/statistic/login_record [post]
func (o *Api) LoginRecord(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetUserLoginRecord, o.adminClient, c)
}

// @Summary		统计用户登录数
// @Description	获取用户登录统计数据
// @Tags			statistic
// @Id	loginUserCount
// @Accept			json
// @Produce		json
// @Param			data	body		chat.UserLoginCountReq	true	"统计条件"
// @Success		200		{object}	apiresp.ApiResponse{data=chat.UserLoginCountResp}
// @Router			/statistic/login_user_count [post]
func (o *Api) LoginUserCount(c *gin.Context) {
	a2r.Call(chat.ChatClient.UserLoginCount, o.chatClient, c)
}

// @Summary		统计新增用户数
// @Description	获取新注册用户统计数据
// @Tags			statistic
// @Id	newUserCount
// @Accept			json
// @Produce		json
// @Param			data	body		user.UserRegisterCountReq	true	"统计条件"
// @Success		200		{object}	apiresp.ApiResponse{data=apistruct.NewUserCountResp}
// @Router			/statistic/new_user_count [post]
func (o *Api) NewUserCount(c *gin.Context) {
	req, err := a2r.ParseRequest[user.UserRegisterCountReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	dateCount, total, err := o.imApiCaller.UserRegisterCount(mctx.WithApiToken(c, imToken), req.Start, req.End)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, &apistruct.NewUserCountResp{
		DateCount: dateCount,
		Total:     total,
	})
}
