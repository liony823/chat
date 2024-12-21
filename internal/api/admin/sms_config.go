package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary		设置短信配置
// @Description	设置短信配置
// @Tags			sms_config
// @Id	setSmsConfig
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SetSmsConfigReq	true	"短信配置信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/sms_config/set [post]
func (o *Api) SetSmsConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.SetSmsConfig, o.adminClient, c)
}

// @Summary		获取短信配置
// @Description	获取短信配置
// @Tags			sms_config
// @Id	getSmsConfig
// @Accept			json
// @Produce		json
// @Success		200		{object}	apiresp.ApiResponse{data=admin.GetSmsConfigResp}
// @Router			/sms_config/get [post]
func (o *Api) GetSmsConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetSmsConfig, o.adminClient, c)
}
