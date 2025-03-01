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

package imapi

import (
	chatpb "github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/msggateway"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/user"
)

// im caller.
var (
	importFriend      = NewApiCaller[relation.ImportFriendReq, relation.ImportFriendResp]("/friend/import_friend")
	getAdminToken     = NewApiCaller[auth.GetAdminTokenReq, auth.GetAdminTokenResp]("/auth/get_admin_token")
	getuserToken      = NewApiCaller[auth.GetUserTokenReq, auth.GetUserTokenResp]("/auth/get_user_token")
	inviteToGroup     = NewApiCaller[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group")
	updateUserInfo    = NewApiCaller[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info")
	updateUserInfoEx  = NewApiCaller[user.UpdateUserInfoExReq, user.UpdateUserInfoExResp]("/user/update_user_info_ex")
	registerUser      = NewApiCaller[user.UserRegisterReq, user.UserRegisterResp]("/user/user_register")
	forceOffLine      = NewApiCaller[auth.ForceLogoutReq, auth.ForceLogoutResp]("/auth/force_logout")
	getGroupsInfo     = NewApiCaller[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/get_groups_info")
	registerUserCount = NewApiCaller[user.UserRegisterCountReq, user.UserRegisterCountResp]("/statistics/user/register")
	friendUserIDs     = NewApiCaller[relation.GetFriendIDsReq, relation.GetFriendIDsResp]("/friend/get_friend_id")
	accountCheck      = NewApiCaller[user.AccountCheckReq, user.AccountCheckResp]("/user/account_check")

	// OWL 新加
	allUserOnlineStatus = NewApiCaller[msggateway.GetUsersOnlineStatusReq, []msggateway.GetUsersOnlineStatusResp_SuccessResult]("/user/get_users_online_status")
	usersOnlineTime     = NewApiCaller[chatpb.GetUsersTimeReq, chatpb.GetUsersTimeResp]("/user/get_users_time")
)
