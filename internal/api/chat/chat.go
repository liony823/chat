// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chat

import (
	"io"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/openimsdk/chat/internal/api/util"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/common/apistruct"
	"github.com/openimsdk/chat/pkg/common/imapi"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	chatpb "github.com/openimsdk/chat/pkg/protocol/chat"
	constantpb "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func New(chatClient chatpb.ChatClient, adminClient admin.AdminClient, imApiCaller imapi.CallerInterface, api *util.Api) *Api {
	return &Api{
		Api:         api,
		chatClient:  chatClient,
		adminClient: adminClient,
		imApiCaller: imApiCaller,
	}
}

type Api struct {
	*util.Api
	chatClient  chatpb.ChatClient
	adminClient admin.AdminClient
	imApiCaller imapi.CallerInterface
}

// ################## ACCOUNT ##################

func (o *Api) SendVerifyCode(c *gin.Context) {
	req, err := a2r.ParseRequest[chatpb.SendVerifyCodeReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.GetClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req.Ip = ip
	resp, err := o.chatClient.SendVerifyCode(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (o *Api) VerifyCode(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.VerifyCode, o.chatClient)
}

func validateSignature(publicKey, nonce, signature string) (b bool, err error) {
	pubKey, err := hexutil.Decode(publicKey)
	if err != nil {
		return false, err
	}
	message, err := hexutil.Decode(nonce)
	if err != nil {
		return false, err
	}
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}
	b = crypto.VerifySignature(pubKey, message, sig)
	return b, nil
}

func (o *Api) RegisterUser(c *gin.Context) {
	req, err := a2r.ParseRequest[chatpb.RegisterUserReq](c)

	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	if validSig, err := validateSignature(req.User.PublicKey, req.Nonce, req.Signature); err != nil || !validSig {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("signature is wrong!"))
		return
	}

	ip, err := o.GetClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req.Ip = ip

	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiCtx := mctx.WithApiToken(c, imToken)
	rpcCtx := o.WithAdminUser(c)

	checkResp, err := o.chatClient.CheckUserExist(rpcCtx, &chatpb.CheckUserExistReq{User: req.User})
	if err != nil {
		log.ZDebug(rpcCtx, "Not else", errs.Unwrap(err))
		apiresp.GinError(c, err)
		return
	}
	if checkResp.IsRegistered {
		isUserNotExist, err := o.imApiCaller.AccountCheckSingle(apiCtx, checkResp.Userid)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		// if User is  not exist in SDK server. You need delete this user and register new user again.
		if isUserNotExist {
			_, err := o.chatClient.DelUserAccount(rpcCtx, &chatpb.DelUserAccountReq{UserIDs: []string{checkResp.Userid}})
			log.ZDebug(c, "Delete Succsssss", checkResp.Userid)
			if err != nil {
				apiresp.GinError(c, err)
				return
			}
		}
	}

	respRegisterUser, err := o.chatClient.RegisterUser(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	userInfo := &sdkws.UserInfo{
		UserID:     respRegisterUser.UserID,
		Account:    respRegisterUser.Account,
		Nickname:   req.User.Nickname,
		FaceURL:    req.User.FaceURL,
		CreateTime: time.Now().UnixMilli(),
		Address:    req.User.Address,
		PublicKey:  req.User.PublicKey,
	}
	err = o.imApiCaller.RegisterUser(apiCtx, []*sdkws.UserInfo{userInfo})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	if resp, err := o.adminClient.FindDefaultFriend(rpcCtx, &admin.FindDefaultFriendReq{}); err == nil {
		_ = o.imApiCaller.ImportFriend(apiCtx, respRegisterUser.UserID, resp.UserIDs)
	}
	if resp, err := o.adminClient.FindDefaultGroup(rpcCtx, &admin.FindDefaultGroupReq{}); err == nil {
		_ = o.imApiCaller.InviteToGroup(apiCtx, respRegisterUser.UserID, resp.GroupIDs)
	}
	var resp apistruct.UserRegisterResp
	if req.AutoLogin {
		resp.ImToken, err = o.imApiCaller.GetUserToken(apiCtx, respRegisterUser.UserID, req.Platform)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
	}
	resp.ChatToken = respRegisterUser.ChatToken
	resp.UserID = respRegisterUser.UserID
	apiresp.GinSuccess(c, &resp)
}

func (o *Api) Login(c *gin.Context) {
	req, err := a2r.ParseRequest[chatpb.LoginReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	ip, err := o.GetClientIP(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req.Ip = ip
	resp, err := o.chatClient.Login(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	adminToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiCtx := mctx.WithApiToken(c, adminToken)

	imToken, err := o.imApiCaller.GetUserToken(apiCtx, resp.UserID, req.Platform)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, &apistruct.LoginResp{
		ImToken:   imToken,
		UserID:    resp.UserID,
		ChatToken: resp.ChatToken,
	})
}

func (o *Api) ResetPassword(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.ResetPassword, o.chatClient)
}

func (o *Api) ChangePassword(c *gin.Context) {
	req, err := a2r.ParseRequest[chatpb.ChangePasswordReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := o.chatClient.ChangePassword(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	err = o.imApiCaller.ForceOffLine(mctx.WithApiToken(c, imToken), req.UserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// ################## USER ##################

func (o *Api) UpdateUserInfo(c *gin.Context) {
	req, err := a2r.ParseRequest[chatpb.UpdateUserInfoReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	respUpdate, err := o.chatClient.UpdateUserInfo(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	var imToken string
	imToken, err = o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	var (
		nickName string
		faceURL  string
	)
	if req.Nickname != nil {
		nickName = req.Nickname.Value
	} else {
		nickName = respUpdate.NickName
	}
	if req.FaceURL != nil {
		faceURL = req.FaceURL.Value
	} else {
		faceURL = respUpdate.FaceUrl
	}
	err = o.imApiCaller.UpdateUserInfo(mctx.WithApiToken(c, imToken), req.UserID, nickName, faceURL)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, apistruct.UpdateUserInfoResp{})
}

func (o *Api) FindUserPublicInfo(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.FindUserPublicInfo, o.chatClient)
}

func (o *Api) FindUserFullInfo(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.FindUserFullInfo, o.chatClient)
}

func (o *Api) SearchUserFullInfo(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.SearchUserFullInfo, o.chatClient)
}

func (o *Api) SearchUserPublicInfo(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.SearchUserPublicInfo, o.chatClient)
}

func (o *Api) GetTokenForVideoMeeting(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.GetTokenForVideoMeeting, o.chatClient)
}

// ################## APPLET ##################

func (o *Api) FindApplet(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.FindApplet, o.adminClient)
}

// ################## CONFIG ##################

func (o *Api) GetClientConfig(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetClientConfig, o.adminClient)
}

// ################## CALLBACK ##################

func (o *Api) OpenIMCallback(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	req := &chatpb.OpenIMCallbackReq{
		Command: c.Query(constantpb.CallbackCommand),
		Body:    string(body),
	}
	if _, err := o.chatClient.OpenIMCallback(c, req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}

func (o *Api) SearchFriend(c *gin.Context) {
	req, err := a2r.ParseRequest[struct {
		UserID string `json:"userID"`
		chatpb.SearchUserInfoReq
	}](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if req.UserID == "" {
		req.UserID = mctx.GetOpUserID(c)
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	userIDs, err := o.imApiCaller.FriendUserIDs(mctx.WithApiToken(c, imToken), req.UserID)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(userIDs) == 0 {
		apiresp.GinSuccess(c, &chatpb.SearchUserInfoResp{})
		return
	}
	req.SearchUserInfoReq.UserIDs = userIDs
	resp, err := o.chatClient.SearchUserInfo(c, &req.SearchUserInfoReq)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (o *Api) LatestApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.LatestApplicationVersion, o.adminClient)
}

func (o *Api) PageApplicationVersion(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.PageApplicationVersion, o.adminClient)
}

/* ################## OWL 新加接口 ################## */

// ################## Post ##################

func (o *Api) PublishPost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.PublishPost, o.chatClient)
}

func (o *Api) ForwardPost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.ForwardPost, o.chatClient)
}

func (o *Api) CommentPost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.CommentPost, o.chatClient)
}

func (o *Api) PinPost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.PinPost, o.chatClient)
}

func (o *Api) ReferencePost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.ReferencePost, o.chatClient)
}

func (o *Api) ChangeLikePost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.ChangeLikePost, o.chatClient)
}

func (o *Api) ChangeCollectPost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.ChangeCollectPost, o.chatClient)
}

func (o *Api) ChangeAllowCommentPost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.ChangeAllowCommentPost, o.chatClient)
}

func (o *Api) ChangeAllowForwardPost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.ChangeAllowForwardPost, o.chatClient)
}

func (o *Api) DeletePost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.DeletePost, o.chatClient)
}

func (o *Api) GetPostByID(c *gin.Context) {
	postID := c.Param("postID")
	if postID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("postID is empty"))
		return
	}
	req := &chatpb.GetPostByIDReq{
		PostID: postID,
	}
	resp, err := o.chatClient.GetPostByID(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (o *Api) GetPostListByUser(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.GetPostListByUser, o.chatClient)
}

func (o *Api) GetPostList(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.GetPostList, o.chatClient)
}

func (o *Api) GetAllTypePost(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.GetAllTypePost, o.chatClient)
}

func (o *Api) GetCommentPostListByPostID(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.GetCommentPostListByPostID, o.chatClient)
}

// ################## Group ##################

func (o *Api) GetGroupFromContact(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.GetGroupFromContact, o.chatClient)
}

func (o *Api) DeleteGroupFromContact(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.DeleteGroupFromContact, o.chatClient)
}

func (o *Api) SaveGroupToContact(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.SaveGroupToContact, o.chatClient)
}

// ################## User ##################

func (o Api) FindUserByAddressOrAccount(c *gin.Context) {
	a2r.Call(c, chatpb.ChatClient.FindUserByAddressOrAccount, o.chatClient)
}

func (o *Api) GetStatistic(c *gin.Context) {
	userResp, err := o.chatClient.GetAllUserIDs(c, &chatpb.GetAllUserIDsReq{})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiCtx := mctx.WithApiToken(c, imToken)
	imResp, err := o.imApiCaller.UserOlineStatus(apiCtx, userResp.UserIDs)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, imResp)
}

func (o *Api) GetUsersOnlineTime(c *gin.Context) {
	req, err := a2r.ParseRequest[chatpb.GetUsersTimeReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	imToken, err := o.imApiCaller.ImAdminTokenWithDefaultAdmin(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiCtx := mctx.WithApiToken(c, imToken)
	imResp, err := o.imApiCaller.UserOlineTimes(apiCtx, req.UserIDs)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, imResp)
}
