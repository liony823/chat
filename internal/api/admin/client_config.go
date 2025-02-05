package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/a2r"
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
	a2r.Call(c, admin.AdminClient.SetClientConfig, o.adminClient)
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
	a2r.Call(c, admin.AdminClient.DelClientConfig, o.adminClient)
}

// @Summary		获取客户端配置
// @Description	获取客户端初始化配置信息
// @Tags			client_config
// @Id	getClientConfig
// @Accept			json
// @Produce		json
// @Success		200		{object}	apiresp.ApiResponse{data=admin.GetClientConfigResp}
// @Router			/client_config/get [post]
func (o *Api) GetClientConfig(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetClientConfig, o.adminClient)
}

// @Summary		获取客户端配置列表
// @Description	获取客户端初始化配置信息列表
// @Tags			client_config
// @Id	getListClientConfig
// @Accept			json
// @Produce		json
// @Success		200		{object}	apiresp.ApiResponse{data=admin.GetListClientConfigResp}
// @Router			/client_config/list [post]
func (o *Api) GetListClientConfig(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetListClientConfig, o.adminClient)
}
