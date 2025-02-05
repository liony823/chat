package admin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
)

// @Summary		添加用户账户
// @Description	创建新的普通用户账户
// @Tags			user
// @Id	addUserAccount
// @Accept			json
// @Produce		json
// @Param			data	body		chat.AddUserAccountReq	true	"用户账户信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/account/add_user [post]
func (o *Api) AddUserAccount(c *gin.Context) {
	req, err := a2r.ParseRequest[chat.AddUserAccountReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.GetClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	ctx := o.WithAdminUser(mctx.WithApiToken(c, imToken))

	err = o.registerChatUser(ctx, ip, []*chat.RegisterUserInfo{req.User})
	if err != nil {
		return
	}

	apiresp.GinSuccess(c, nil)
}

// @Summary		重置用户密码
// @Description	重置指定用户的密码
// @Tags			user
// @Id	resetUserPassword
// @Accept			json
// @Produce		json
// @Param			data	body		chat.ChangePasswordReq	true	"密码重置信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/user/password/reset [post]
func (o *Api) ResetUserPassword(c *gin.Context) {
	req, err := a2r.ParseRequest[chat.ChangePasswordReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := o.chatClient.ChangePassword(c, req)
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

// @Summary		删除管理员账户
// @Description	删除指定管理员账户
// @Tags			user
// @Id	delAdminAccount
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DelAdminAccountReq	true	"删除管理员账户"
// @Success		200	{object}	apiresp.ApiResponse
// @Router			/account/del_admin [post]
func (o *Api) DelAdminAccount(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelAdminAccount, o.adminClient)
}

// @Summary		搜索管理员账户
// @Description	搜索管理员账户
// @Tags			user
// @Id	searchAdminAccount
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SearchAdminAccountReq	true	"管理员账户搜索"
// @Success		200	{object}	apiresp.ApiResponse{data=admin.SearchAdminAccountResp}
// @Router			/account/search_admin [post]
func (o *Api) SearchAdminAccount(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchAdminAccount, o.adminClient)
}

func (o *Api) registerChatUser(ctx context.Context, ip string, users []*chat.RegisterUserInfo) error {
	if len(users) == 0 {
		return errs.ErrArgs.WrapMsg("users is empty")
	}
	for _, info := range users {
		respRegisterUser, err := o.chatClient.RegisterUser(ctx, &chat.RegisterUserReq{Ip: ip, User: info, Platform: constant.AdminPlatformID})
		if err != nil {
			return err
		}
		userInfo := &sdkws.UserInfo{
			UserID:   respRegisterUser.UserID,
			Nickname: info.Nickname,
			FaceURL:  info.FaceURL,
		}
		if err = o.imApiCaller.RegisterUser(ctx, []*sdkws.UserInfo{userInfo}); err != nil {
			return err
		}

		if resp, err := o.adminClient.FindDefaultFriend(ctx, &admin.FindDefaultFriendReq{}); err == nil {
			_ = o.imApiCaller.ImportFriend(ctx, respRegisterUser.UserID, resp.UserIDs)
		}
		if resp, err := o.adminClient.FindDefaultGroup(ctx, &admin.FindDefaultGroupReq{}); err == nil {
			_ = o.imApiCaller.InviteToGroup(ctx, respRegisterUser.UserID, resp.GroupIDs)
		}
	}
	return nil
}
