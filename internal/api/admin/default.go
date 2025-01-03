package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/liony823/tools/a2r"
	"github.com/liony823/tools/apiresp"
	"github.com/liony823/tools/errs"
	"github.com/openimsdk/chat/pkg/common/apistruct"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/protocol/sdkws"
)

// @Summary		添加默认好友
// @Description	添加用户注册时的默认好友
// @Tags			default
// @Id	addDefaultFriend
// @Accept			json
// @Produce		json
// @Param			data	body		admin.AddDefaultFriendReq	true	"默认好友信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/default/user/add [post]
func (o *Api) AddDefaultFriend(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddDefaultFriend, o.adminClient)
}

// @Summary		添加默认群组
// @Description	添加用户注册时自动加入的默认群组
// @Tags			default
// @Id	addDefaultGroup
// @Accept			json
// @Produce		json
// @Param			data	body		admin.AddDefaultGroupReq	true	"默认群组信息"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/default/group/add [post]
func (o *Api) AddDefaultGroup(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.AddDefaultGroupReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	groups, err := o.imApiCaller.FindGroupInfo(mctx.WithApiToken(c, imToken), req.GroupIDs)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(req.GroupIDs) != len(groups) {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("group id not found"))
		return
	}
	resp, err := o.adminClient.AddDefaultGroup(c, &admin.AddDefaultGroupReq{
		GroupIDs: req.GroupIDs,
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// @Summary		删除默认好友
// @Description	删除用户注册时的默认好友
// @Tags			default
// @Id	delDefaultFriend
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DelDefaultFriendReq	true	"要删除的好友ID"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/default/user/del [post]
func (o *Api) DelDefaultFriend(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelDefaultFriend, o.adminClient)
}

// @Summary		查找默认好友
// @Description	获取所有默认好友列表
// @Tags			default
// @Id	findDefaultFriend
// @Accept			json
// @Produce		json
// @Success		200	{object}	apiresp.ApiResponse{data=admin.FindDefaultFriendResp}
// @Router			/default/user/find [post]
func (o *Api) FindDefaultFriend(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.FindDefaultFriend, o.adminClient)
}

// @Summary		搜索默认好友
// @Description	搜索用户注册时的默认好友
// @Tags			default
// @Id	searchDefaultFriend
// @Accept			json
// @Produce		json
// @Success		200	{object}	apiresp.ApiResponse{data=admin.SearchDefaultFriendResp}
// @Router			/default/user/search [post]
func (o *Api) SearchDefaultFriend(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.SearchDefaultFriend, o.adminClient)
}

// @Summary		删除默认群组
// @Description	删除用户注册时自动加入的默认群组
// @Tags			default
// @Id	delDefaultGroup
// @Accept			json
// @Produce		json
// @Param			data	body		admin.DelDefaultGroupReq	true	"要删除的群组ID"
// @Success		200		{object}	apiresp.ApiResponse
// @Router			/default/group/del [post]
func (o *Api) DelDefaultGroup(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DelDefaultGroup, o.adminClient)
}

// @Summary		查找默认群组
// @Description	获取所有默认群组列表
// @Tags			default
// @Id	findDefaultGroup
// @Accept			json
// @Produce		json
// @Success		200	{object}	apiresp.ApiResponse{data=admin.FindDefaultGroupResp}
// @Router			/default/group/find [post]
func (o *Api) FindDefaultGroup(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.FindDefaultGroup, o.adminClient)
}

// @Summary		搜索默认群组
// @Description	搜索用户注册时自动加入的默认群组
// @Tags			default
// @Id	searchDefaultGroup
// @Accept			json
// @Produce		json
// @Success		200	{object}	apiresp.ApiResponse{data=admin.SearchDefaultGroupResp}
// @Router			/default/group/search [post]
func (o *Api) SearchDefaultGroup(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.SearchDefaultGroupReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	searchResp, err := o.adminClient.SearchDefaultGroup(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp := apistruct.SearchDefaultGroupResp{
		Total:  searchResp.Total,
		Groups: make([]*sdkws.GroupInfo, 0, len(searchResp.GroupIDs)),
	}
	if len(searchResp.GroupIDs) > 0 {
		imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		groups, err := o.imApiCaller.FindGroupInfo(mctx.WithApiToken(c, imToken), searchResp.GroupIDs)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		groupMap := make(map[string]*sdkws.GroupInfo)
		for _, group := range groups {
			groupMap[group.GroupID] = group
		}
		for _, groupID := range searchResp.GroupIDs {
			if group, ok := groupMap[groupID]; ok {
				resp.Groups = append(resp.Groups, group)
			} else {
				resp.Groups = append(resp.Groups, &sdkws.GroupInfo{
					GroupID: groupID,
				})
			}
		}
	}
	apiresp.GinSuccess(c, resp)
}
