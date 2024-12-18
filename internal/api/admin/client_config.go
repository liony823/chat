package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary		设置客户端配置
// @Description	设置客户端初始化配置信息
// @Tags			client_config
// @Id	setClientConfig
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SetClientConfigReq	true	"客户端配置信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/client_config/set [post]
func (o *Api) SetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.SetClientConfig, o.adminClient, c)
}

// @Summary		删除客户端配置
// @Description	删除指定的客户端配置
// @Tags			client_config
// @Id	delClientConfig
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DelClientConfigReq	true	"要删除的配置信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/client_config/del [post]
func (o *Api) DelClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.DelClientConfig, o.adminClient, c)
}

// @Summary		获取客户端配置
// @Description	获取客户端初始化配置信息
// @Tags			client_config
// @Id	getClientConfig
// @Accept			json
// @Produce		json
// @Param			data	body		admin.GetClientConfigReq	true	"查询条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.GetClientConfigResp}
// @Router			/client_config/get [post]
func (o *Api) GetClientConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetClientConfig, o.adminClient, c)
}
