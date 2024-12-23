package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary		设置签到配置
// @Description	设置签到配置
// @Tags			signin_config
// @Id				setSigninConfig
// @Accept			json
// @Produce			json
// @Param			data	body		admin.SetSigninConfigReq	true	"签到配置信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/signin_config/set [post]
func (o *Api) SetSigninConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.SetSigninConfig, o.adminClient, c)
}

// @Summary		获取签到配置
// @Description	获取签到配置
// @Tags			signin_config
// @Id				getSigninConfig
// @Accept			json
// @Produce			json
// @Success		200		{object}	apiresp.ApiResponse{data=admin.GetSigninConfigResp}
// @Router			/signin_config/get [post]
func (o *Api) GetSigninConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetSigninConfig, o.adminClient, c)
}
