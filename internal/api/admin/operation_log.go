package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/a2r"
)

// @Summary		搜索操作日志
// @Description	搜索操作日志
// @Tags			operation_log
// @Id	searchOperationLog
// @Accept			json
// @Produce		json
// @Param			data	body		admin.SearchOperationLogReq	true	"查询条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.SearchOperationLogResp}
// @Router			/operation_log/search [post]
func (o *Api) SearchOperationLog(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchOperationLog, o.adminClient)
}

// @Summary		获取操作日志
// @Description	获取操作日志
// @Tags			operation_log
// @Id	getOperationLog
// @Accept			json
// @Produce		json
// @Param			data	body		admin.GetOperationLogReq	true	"查询条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.GetOperationLogResp}
// @Router			/operation_log/get [post]
func (o *Api) GetOperationLog(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetOperationLog, o.adminClient)
}

// @Summary		删除操作日志
// @Description	删除操作日志
// @Tags			operation_log
// @Id	delOperationLog
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DeleteOperationLogReq	true	"查询条件"
// @Success		200		{object}	apiresp.ApiResponse{data=admin.DeleteOperationLogResp}
// @Router			/operation_log/delete [post]
func (o *Api) DeleteOperationLog(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DeleteOperationLog, o.adminClient)
}
