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
	"context"
	"errors"
	"time"

	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/tools/utils/stringutil"

	"github.com/openimsdk/chat/pkg/common/convert"
	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	chatdb "github.com/openimsdk/chat/pkg/common/db/table/chat"
	constantpb "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/chat"
	"github.com/openimsdk/chat/pkg/protocol/common"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
)

func (o *chatSvr) checkUpdateInfo(ctx context.Context, req *chat.UpdateUserInfoReq) error {
	if req.UserID == "" {
		return errs.ErrArgs.WrapMsg("user id is empty")
	}

	credentials, err := o.Database.TakeCredentialsByUserID(ctx, req.UserID)
	if err != nil {
		return err
	} else if len(credentials) == 0 {
		return errs.ErrArgs.WrapMsg("user not found")
	}

	// 计算凭证数量变化
	credNum := len(credentials)
	delNum := 0
	addNum := 0

	// 检查新增账号
	if req.Account != nil && req.Account.Value != "" {
		addNum++
	}

	// 检查现有凭证
	for _, credential := range credentials {
		if credential.Type == constant.CredentialAccount && req.Account != nil {
			if req.Account.Value == credential.Account {
				req.Account = nil // 相同账号不需要更新
			} else if req.Account.Value == "" {
				delNum++ // 删除现有账号
			}
		}
	}

	// 确保至少保留一种登录方式
	if addNum+credNum-delNum <= 0 {
		return errs.ErrArgs.WrapMsg("a login method must exist")
	}

	// 验证新账号
	if req.Account != nil && req.Account.GetValue() != "" {
		// 检查账号格式
		if !stringutil.IsAlphanumeric(req.Account.GetValue()) {
			return errs.ErrArgs.WrapMsg("account must be alphanumeric")
		}

		// 检查账号是否已存在
		_, err := o.Database.TakeCredentialByAccount(ctx, req.Account.GetValue())
		if err == nil {
			return eerrs.ErrAccountAlreadyRegister.Wrap()
		} else if !dbutil.IsDBNotFound(err) {
			return err
		}
	}

	return nil
}

func (o *chatSvr) UpdateUserInfo(ctx context.Context, req *chat.UpdateUserInfoReq) (*chat.UpdateUserInfoResp, error) {
	opUserID, userType, err := mctx.Check(ctx)
	if err != nil {
		return nil, err
	}

	if err = o.checkUpdateInfo(ctx, req); err != nil {
		return nil, err
	}

	switch userType {
	case constant.NormalUser:
		if req.RegisterType != nil {
			return nil, errs.ErrNoPermission.WrapMsg("registerType can not be updated")
		}
		if req.UserID != opUserID {
			return nil, errs.ErrNoPermission.WrapMsg("only admin can update other user info")
		}

	case constant.AdminUser:
	default:
		return nil, errs.ErrNoPermission.WrapMsg("user type error")
	}

	update, err := ToDBAttributeUpdate(req)
	if err != nil {
		return nil, err
	}
	credUpdate, credDel, err := ToDBCredentialUpdate(req, true)
	if err != nil {
		return nil, err
	}
	if len(update) > 0 {
		if err := o.Database.UpdateUseInfo(ctx, req.UserID, update, credUpdate, credDel); err != nil {
			return nil, err
		}
	}
	return &chat.UpdateUserInfoResp{}, nil
}

func (o *chatSvr) FindUserPublicInfo(ctx context.Context, req *chat.FindUserPublicInfoReq) (*chat.FindUserPublicInfoResp, error) {
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("UserIDs is empty")
	}
	attributes, err := o.Database.FindAttribute(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &chat.FindUserPublicInfoResp{
		Users: convert.DbToPbAttributes(attributes),
	}, nil
}

func (o *chatSvr) AddUserAccount(ctx context.Context, req *chat.AddUserAccountReq) (*chat.AddUserAccountResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}

	if err := o.checkRegisterInfo(ctx, req.User, true); err != nil {
		return nil, err
	}

	if req.User.UserID == "" {
		for i := 0; i < 20; i++ {
			userID := o.genUserID()
			_, err := o.Database.GetUser(ctx, userID)
			if err == nil {
				continue
			} else if dbutil.IsDBNotFound(err) {
				req.User.UserID = userID
				break
			} else {
				return nil, err
			}
		}
		if req.User.UserID == "" {
			return nil, errs.ErrInternalServer.WrapMsg("gen user id failed")
		}
	} else {
		_, err := o.Database.GetUser(ctx, req.User.UserID)
		if err == nil {
			return nil, errs.ErrArgs.WrapMsg("appoint user id already register")
		} else if !dbutil.IsDBNotFound(err) {
			return nil, err
		}
	}

	var (
		credentials []*chatdb.Credential
	)

	if req.User.Account != "" {
		credentials = append(credentials, &chatdb.Credential{
			UserID:      req.User.UserID,
			Account:     req.User.Account,
			Type:        constant.CredentialAccount,
			AllowChange: true,
		})
	}

	register := &chatdb.Register{
		UserID:      req.User.UserID,
		DeviceID:    req.DeviceID,
		IP:          req.Ip,
		Platform:    constantpb.PlatformID2Name[int(req.Platform)],
		AccountType: "",
		Mode:        constant.UserMode,
		CreateTime:  time.Now(),
	}
	account := &chatdb.Account{
		UserID:         req.User.UserID,
		OperatorUserID: mcontext.GetOpUserID(ctx),
		ChangeTime:     register.CreateTime,
		CreateTime:     register.CreateTime,
	}
	attribute := &chatdb.Attribute{
		UserID:            req.User.UserID,
		Account:           req.User.Account,
		Nickname:          req.User.Nickname,
		FaceURL:           req.User.FaceURL,
		ChangeTime:        register.CreateTime,
		CreateTime:        register.CreateTime,
		ChangeAccountTime: register.CreateTime,
		AllowVibration:    constant.DefaultAllowVibration,
		AllowBeep:         constant.DefaultAllowBeep,
		AllowAddFriend:    constant.DefaultAllowAddFriend,
	}

	if err := o.Database.RegisterUser(ctx, register, account, attribute, credentials); err != nil {
		return nil, err
	}
	return &chat.AddUserAccountResp{}, nil
}

func (o *chatSvr) SearchUserPublicInfo(ctx context.Context, req *chat.SearchUserPublicInfoReq) (*chat.SearchUserPublicInfoResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	total, list, err := o.Database.Search(ctx, constant.FinDAllUser, req.Keyword, req.Genders, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &chat.SearchUserPublicInfoResp{
		Total: uint32(total),
		Users: convert.DbToPbAttributes(list),
	}, nil
}

func (o *chatSvr) FindUserFullInfo(ctx context.Context, req *chat.FindUserFullInfoReq) (*chat.FindUserFullInfoResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("UserIDs is empty")
	}
	attributes, err := o.Database.FindAttribute(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &chat.FindUserFullInfoResp{Users: convert.DbToPbUserFullInfos(attributes)}, nil
}

func (o *chatSvr) SearchUserFullInfo(ctx context.Context, req *chat.SearchUserFullInfoReq) (*chat.SearchUserFullInfoResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	total, list, err := o.Database.Search(ctx, req.Normal, req.Keyword, req.Genders, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &chat.SearchUserFullInfoResp{
		Total: uint32(total),
		Users: convert.DbToPbUserFullInfos(list),
	}, nil
}

func (o *chatSvr) FindUserAccount(ctx context.Context, req *chat.FindUserAccountReq) (*chat.FindUserAccountResp, error) {
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("user id list must be set")
	}
	if _, _, err := mctx.CheckAdminOrUser(ctx); err != nil {
		return nil, err
	}
	attributes, err := o.Database.FindAttribute(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	userAccountMap := make(map[string]string)
	for _, attribute := range attributes {
		userAccountMap[attribute.UserID] = attribute.Account
	}
	return &chat.FindUserAccountResp{UserAccountMap: userAccountMap}, nil
}

func (o *chatSvr) FindAccountUser(ctx context.Context, req *chat.FindAccountUserReq) (*chat.FindAccountUserResp, error) {
	if len(req.Accounts) == 0 {
		return nil, errs.ErrArgs.WrapMsg("account list must be set")
	}
	if _, _, err := mctx.CheckAdminOrUser(ctx); err != nil {
		return nil, err
	}
	attributes, err := o.Database.FindAttribute(ctx, req.Accounts)
	if err != nil {
		return nil, err
	}
	accountUserMap := make(map[string]string)
	for _, attribute := range attributes {
		accountUserMap[attribute.Account] = attribute.UserID
	}
	return &chat.FindAccountUserResp{AccountUserMap: accountUserMap}, nil
}

func (o *chatSvr) SearchUserInfo(ctx context.Context, req *chat.SearchUserInfoReq) (*chat.SearchUserInfoResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}
	total, list, err := o.Database.SearchUser(ctx, req.Keyword, req.UserIDs, req.Genders, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &chat.SearchUserInfoResp{
		Total: uint32(total),
		Users: convert.DbToPbUserFullInfos(list),
	}, nil
}

func (o *chatSvr) CheckUserExist(ctx context.Context, req *chat.CheckUserExistReq) (resp *chat.CheckUserExistResp, err error) {
	if req.User == nil {
		return nil, errs.ErrArgs.WrapMsg("user is nil")
	}
	// if req.User.PhoneNumber != "" {
	// 	account, err := o.Database.TakeCredentialByAccount(ctx, BuildCredentialPhone(req.User.AreaCode, req.User.PhoneNumber))
	// 	// err != nil is not found User
	// 	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
	// 		return nil, err
	// 	}
	// 	if account != nil {
	// 		log.ZDebug(ctx, "Check Number is ", account.Account)
	// 		log.ZDebug(ctx, "Check userID is ", account.UserID)
	// 		return &chat.CheckUserExistResp{Userid: account.UserID, IsRegistered: true}, nil
	// 	}
	// }
	// if req.User.Email != "" {
	// 	account, err := o.Database.TakeCredentialByAccount(ctx, req.User.AreaCode)
	// 	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
	// 		return nil, err
	// 	}
	// 	if account != nil {
	// 		log.ZDebug(ctx, "Check email is ", account.Account)
	// 		log.ZDebug(ctx, "Check userID is ", account.UserID)
	// 		return &chat.CheckUserExistResp{Userid: account.UserID, IsRegistered: true}, nil
	// 	}
	// }
	if req.User.Account != "" {
		account, err := o.Database.TakeCredentialByAccount(ctx, req.User.Account)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		if account != nil {
			log.ZDebug(ctx, "Check account is ", account.Account)
			log.ZDebug(ctx, "Check userID is ", account.UserID)
			return &chat.CheckUserExistResp{Userid: account.UserID, IsRegistered: true}, nil
		}
	}
	return nil, nil
}

func (o *chatSvr) DelUserAccount(ctx context.Context, req *chat.DelUserAccountReq) (resp *chat.DelUserAccountResp, err error) {
	if err := o.Database.DelUserAccount(ctx, req.UserIDs); err != nil && errs.Unwrap(err) != mongo.ErrNoDocuments {
		return nil, err
	}
	return nil, nil
}

func (o *chatSvr) FindUserByAddressOrAccount(ctx context.Context, req *chat.FindUserByAddressOrAccountReq) (*chat.FindUserPublicInfoResp, error) {
	if _, _, err := mctx.Check(ctx); err != nil {
		return nil, err
	}

	if req.Account == "" && req.Address == "" {
		return nil, errs.ErrArgs.WrapMsg("search is empty")
	}

	attribute := &chatdb.Attribute{}
	var err error

	if req.Account != "" {
		attribute, err = o.Database.TakeAttributeByAccount(ctx, req.Account)
	}

	if attribute == nil {
		if req.Address != "" {
			attribute, err = o.Database.TakeAttributeByAddress(ctx, req.Address)
		}
	}

	if attribute == nil {
		return &chat.FindUserPublicInfoResp{
			Users: []*common.UserPublicInfo{},
		}, nil
	}

	if err != nil {
		return nil, err
	}
	if req.Account != attribute.Account || req.Address != attribute.Address {
		return &chat.FindUserPublicInfoResp{
			Users: []*common.UserPublicInfo{convert.DbToPbAttribute(attribute)},
		}, nil
	}
	return nil, nil
}

func (o *chatSvr) GetAllUserIDs(ctx context.Context, req *chat.GetAllUserIDsReq) (*chat.GetAllUserIDsResp, error) {
	var userIDs []string
	pageNumber := 1
	showNumber := 500
	for {
		_, recvIDsPart, err := o.Database.GetAllUserID(ctx, &sdkws.RequestPagination{PageNumber: int32(pageNumber), ShowNumber: int32(showNumber)})
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, recvIDsPart...)
		if len(recvIDsPart) < showNumber {
			break
		}
		pageNumber++
	}

	return &chat.GetAllUserIDsResp{UserIDs: userIDs}, nil
}


func (o *chatSvr) SetStealthUser(ctx context.Context, req *chat.SetStealthUserReq) (*chat.SetStealthUserResp, error) {
	if _, _, err := mctx.CheckAdminOrUser(ctx); err != nil {
		return nil, err
	}

	if err := o.Database.SetStealthUser(ctx, req.UserID, int64(req.Stealth)); err != nil {
		return nil, err
	}
	return &chat.SetStealthUserResp{}, nil
}
