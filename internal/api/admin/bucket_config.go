package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// @Summary		获取桶配置
// @Description	获取桶配置
// @Tags			bucket_config
// @Id				getBucketConfig
// @Accept			json
// @Produce		json
// @Success		200		{object}	apiresp.ApiResponse{data=admin.GetBucketConfigResp}
// @Router			/bucket_config/get [post]
func (o *Api) GetBucketConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.GetBucketConfig, o.adminClient, c)
}


// @Summary		设置桶配置
// @Description	设置桶配置
// @Tags			bucket_config
// @Id				setBucketConfig
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SetBucketConfigReq	true	"桶配置信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/bucket_config/set [post]
func (o *Api) SetBucketConfig(c *gin.Context) {
	a2r.Call(admin.AdminClient.SetBucketConfig, o.adminClient, c)
}
