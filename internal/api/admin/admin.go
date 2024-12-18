package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/apiresp"
	_ "github.com/liony823/tools/apiresp"
	"github.com/openimsdk/chat/internal/api/util"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/chat/pkg/protocol/chat"
)

type Api struct {
	*util.Api
	chatClient  chat.ChatClient
	adminClient admin.AdminClient
	imApiCaller imapi.CallerInterface
}

func New(chatClient chat.ChatClient, adminClient admin.AdminClient, imApiCaller imapi.CallerInterface, api *util.Api) *Api {
	return &Api{
		Api:         api,
		chatClient:  chatClient,
		adminClient: adminClient,
		imApiCaller: imApiCaller,
	}
}

// @Summary		测试
// @Description	测试
// @Tags			common
// @Id	ping
// @Success		200	{object} apiresp.ApiResponse{data=string}
// @Router			/ping [get]
func (o *Api) AdminPing(c *gin.Context) {
	apiresp.GinSuccess(c, "pong")
}
