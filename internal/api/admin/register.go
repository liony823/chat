package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/tools/a2r"
)

// @Summary		设置是否允许注册
// @Description	设置系统是否允许新用户注册
// @Tags			user
// @Id	setAllowRegister
// @Accept			json
// @Produce		json
// @Param			data	body		chat.SetAllowRegisterReq	true	"注册配置"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/user/allow_register/set [post]
func (o *Api) SetAllowRegister(c *gin.Context) {
	a2r.Call(c, chat.ChatClient.SetAllowRegister, o.chatClient)
}

// @Summary		获取注册配置
// @Description	获取系统是否允许新用户注册的配置
// @Tags			user
// @Id	getAllowRegister
// @Accept			json
// @Produce		json
// @Success		200	{object}	apiresp.ApiResponse{data=chat.GetAllowRegisterResp}
// @Router			/user/allow_register/get [post]
func (o *Api) GetAllowRegister(c *gin.Context) {
	a2r.Call(c, chat.ChatClient.GetAllowRegister, o.chatClient)
}
