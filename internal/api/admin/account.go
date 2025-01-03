package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/liony823/tools/apiresp"
	"github.com/liony823/tools/log"
	"github.com/liony823/tools/utils/datautil"
	"github.com/openimsdk/chat/pkg/common/apistruct"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary		管理员登录
// @Description	管理员账户登录
// @Tags			account
// @Id	login
// @Accept			json
// @Produce		json
// @Param			data	body		admin.LoginReq	true	"登录信息"
// @Success		200		{object}	apiresp.ApiResponse{data=apistruct.AdminLoginResp}
// @Router			/account/login [post]
func (o *Api) AdminLogin(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.LoginReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	loginResp, err := o.adminClient.Login(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imAdminUserID := o.GetDefaultIMAdminUserID()
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	var resp apistruct.AdminLoginResp
	if err := datautil.CopyStructFields(&resp, loginResp); err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp.ImToken = imToken
	resp.ImUserID = imAdminUserID
	apiresp.GinSuccess(c, resp)

}

// @Summary		更新管理员信息
// @Description	更新管理员账户信息
// @Tags			account
// @Id	updateInfo
// @Accept			json
// @Produce		json
// @Param			data	body		admin.AdminUpdateInfoReq	true	"管理员信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/account/update [post]
func (o *Api) AdminUpdateInfo(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.AdminUpdateInfoReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := o.adminClient.AdminUpdateInfo(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	imAdminUserID := o.GetDefaultIMAdminUserID()
	imToken, err := o.imApiCaller.GetAdminTokenCache(c, imAdminUserID)
	if err != nil {
		log.ZError(c, "AdminUpdateInfo ImAdminTokenWithDefaultAdmin", err, "imAdminUserID", imAdminUserID)
		return
	}
	if err := o.imApiCaller.UpdateUserInfo(mctx.WithApiToken(c, imToken), imAdminUserID, resp.Nickname, resp.FaceURL); err != nil {
		log.ZError(c, "AdminUpdateInfo UpdateUserInfo", err, "userID", resp.UserID, "nickName", resp.Nickname, "faceURL", resp.FaceURL)
	}
	apiresp.GinSuccess(c, nil)
}

// @Summary		管理员信息
// @Description	获取管理员信息
// @Tags			account
// @Id	getInfo
// @Accept			json
// @Produce		json
// @Success		200	{object}	apiresp.ApiResponse{data=admin.GetAdminInfoResp}
// @Router			/account/info [post]
func (o *Api) AdminInfo(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetAdminInfo, o.adminClient)
}

// @Summary		管理员密码
// @Description	修改管理员密码
// @Tags			account
// @Id	changePassword
// @Accept			json
// @Produce		json
// @Param			data	body		admin.ChangeAdminPasswordReq	true	"管理员密码"
// @Success		200	{object}	apiresp.ApiResponse
// @Router			/account/change_password [post]
func (o *Api) ChangeAdminPassword(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.ChangeAdminPassword, o.adminClient)
}

// @Summary		添加管理员账户
// @Description	创建新的管理员账户
// @Tags			account
// @Id	addAdminAccount
// @Accept			json
// @Produce		json
// @Param			data	body		admin.AddAdminAccountReq	true	"添加管理员账户信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/account/add_admin [post]
func (o *Api) AddAdminAccount(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddAdminAccount, o.adminClient)
}

// @Summary		获取Google Authenticator密钥和二维码数据
// @Description	获取管理员账户的Google Authenticator密钥和二维码数据
// @Tags			account
// @Id	getGoogleAuth
// @Accept			json
// @Produce		json
// @Success		200	{object}	apiresp.ApiResponse{data=admin.GetGoogleAuthResp}
// @Router			/account/google_auth [post]
func (o *Api) GetGoogleAuth(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetGoogleAuth, o.adminClient)
}

// @Summary		验证Google Authenticator动态令牌
// @Description	验证管理员账户的Google Authenticator动态令牌,开启两步验证
// @Tags			account
// @Id	verifyGoogleAuth
// @Accept			json
// @Produce		json
// @Param			data	body		admin.VerifyGoogleAuthReq	true	"验证信息"
// @Success		200	{object}	apiresp.ApiResponse{data=admin.VerifyGoogleAuthResp}
// @Router			/account/verify_google_auth [post]
func (o *Api) VerifyGoogleAuth(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.VerifyGoogleAuth, o.adminClient)
}
