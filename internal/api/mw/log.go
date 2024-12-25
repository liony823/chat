package mw

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/common/utils"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

// OperationLog 操作日志中间件
func (m *MW) OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不记录GET请求和swagger相关请求
		if c.Request.Method == "GET" || strings.Contains(c.Request.URL.Path, "swagger") {
			c.Next()
			return
		}

		// 获取请求数据
		var requestData string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			requestData = string(bodyBytes)
		}

		// 获取管理员信息
		adminInfo, err := m.client.GetAdminInfo(c, &admin.GetAdminInfoReq{})
		if err != nil {
			c.Next()
			return
		}
		if adminInfo == nil {
			c.Next()
			return
		}

		// 获取模块和操作说明
		module, operation := utils.GetModuleAndOperation(c.Request.URL.Path)

		// 异步保存日志
		go func() {
			m.client.CreateOperationLog(c, &admin.CreateOperationLogReq{
				OperationID:  adminInfo.UserID,
				AdminID:      adminInfo.UserID,
				AdminAccount: adminInfo.Account,
				AdminName:    adminInfo.Nickname,
				Module:       module,
				Operation:    operation,
				Method:       c.Request.Method,
				Path:         c.Request.URL.Path,
				IP:           c.ClientIP(),
				RequestData:  requestData,
				CreateTime:   time.Now().UnixMilli(),
			})
		}()

		c.Next()
	}
}
