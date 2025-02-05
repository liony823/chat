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

package admin

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/pwdutil"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) GetAdminInfo(ctx context.Context, req *admin.GetAdminInfoReq) (*admin.GetAdminInfoResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	a, err := o.Database.GetAdminUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res, err := o.ListAdminUserMenu(ctx, &admin.ListAdminUserMenuReq{UserID: userID})
	if err != nil {
		return nil, err
	}
	return &admin.GetAdminInfoResp{
		Account:          a.Account,
		FaceURL:          a.FaceURL,
		Nickname:         a.Nickname,
		UserID:           a.UserID,
		EnableGoogleAuth: a.EnableGoogleAuth,
		Level:            a.Level,
		CreateTime:       a.CreateTime.UnixMilli(),
		Menus:            res.Menus,
	}, nil
}

func (o *adminServer) ChangeAdminPassword(ctx context.Context, req *admin.ChangeAdminPasswordReq) (*admin.ChangeAdminPasswordResp, error) {
	user, err := o.Database.GetAdminUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if !pwdutil.VerifyPassword(req.CurrentPassword, user.Password) {
		return nil, eerrs.ErrPassword.Wrap()
	}
	hashedPassword, err := pwdutil.EncryptPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}
	if err := o.Database.ChangePassword(ctx, req.UserID, hashedPassword); err != nil {
		return nil, err
	}
	return &admin.ChangeAdminPasswordResp{}, nil
}

func (o *adminServer) AddAdminAccount(ctx context.Context, req *admin.AddAdminAccountReq) (*admin.AddAdminAccountResp, error) {
	if err := o.CheckSuperAdmin(ctx); err != nil {
		return nil, err
	}

	_, err := o.Database.GetAdmin(ctx, req.Account)
	if err == nil {
		return nil, errs.ErrDuplicateKey.WrapMsg("the account is registered")
	}

	hashedPassword, err := pwdutil.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	adm := &admindb.Admin{
		Account:          req.Account,
		Password:         hashedPassword,
		FaceURL:          req.FaceURL,
		Nickname:         req.Nickname,
		UserID:           o.genUserID(),
		Level:            constant.NormalAdmin,
		CreateTime:       time.Now(),
		EnableGoogleAuth: false,
		GoogleAuthSecret: "",
	}
	if err = o.Database.AddAdminAccount(ctx, []*admindb.Admin{adm}); err != nil {
		return nil, err
	}
	return &admin.AddAdminAccountResp{}, nil
}

func (o *adminServer) DelAdminAccount(ctx context.Context, req *admin.DelAdminAccountReq) (*admin.DelAdminAccountResp, error) {
	if err := o.CheckSuperAdmin(ctx); err != nil {
		return nil, err
	}

	if datautil.Duplicate(req.UserIDs) {
		return nil, errs.ErrArgs.WrapMsg("user ids is duplicate")
	}

	for _, userID := range req.UserIDs {
		superAdmin, err := o.Database.GetAdminUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if superAdmin.Level == constant.AdvancedUserLevel {
			return nil, errs.ErrNoPermission.WrapMsg(fmt.Sprintf("%s is superAdminID", userID))
		}
	}

	if err := o.Database.DelAdminAccount(ctx, req.UserIDs); err != nil {
		return nil, err
	}
	return &admin.DelAdminAccountResp{}, nil
}

func (o *adminServer) SearchAdminAccount(ctx context.Context, req *admin.SearchAdminAccountReq) (*admin.SearchAdminAccountResp, error) {
	if err := o.CheckSuperAdmin(ctx); err != nil {
		return nil, err
	}

	filter := bson.M{}
	if req.Account != "" {
		filter["account"] = req.Account
	}
	if req.Nickname != "" {
		filter["nickname"] = req.Nickname
	}

	total, adminAccounts, err := o.Database.SearchAdminAccount(ctx, req.Pagination, filter)
	if err != nil {
		return nil, err
	}
	accounts := make([]*admin.GetAdminInfoResp, 0, len(adminAccounts))
	for _, v := range adminAccounts {
		temp := &admin.GetAdminInfoResp{
			Account:  v.Account,
			FaceURL:  v.FaceURL,
			Nickname: v.Nickname,
			UserID:   v.UserID,
			Level:    v.Level,

			CreateTime: v.CreateTime.Unix(),
		}
		accounts = append(accounts, temp)
	}
	return &admin.SearchAdminAccountResp{Total: uint32(total), AdminAccounts: accounts}, nil
}

func (o *adminServer) AdminUpdateInfo(ctx context.Context, req *admin.AdminUpdateInfoReq) (*admin.AdminUpdateInfoResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	update, err := ToDBAdminUpdate(req)
	if err != nil {
		return nil, err
	}
	info, err := o.Database.GetAdminUserID(ctx, mcontext.GetOpUserID(ctx))
	if err != nil {
		return nil, err
	}
	if err := o.Database.UpdateAdmin(ctx, userID, update); err != nil {
		return nil, err
	}
	resp := &admin.AdminUpdateInfoResp{UserID: info.UserID}
	if req.Nickname == nil {
		resp.Nickname = info.Nickname
	} else {
		resp.Nickname = req.Nickname.Value
	}
	if req.FaceURL == nil {
		resp.FaceURL = info.FaceURL
	} else {
		resp.FaceURL = req.FaceURL.Value
	}
	return resp, nil
}

func (o *adminServer) Login(ctx context.Context, req *admin.LoginReq) (*admin.LoginResp, error) {
	a, err := o.Database.GetAdmin(ctx, req.Account)
	if err != nil {
		if dbutil.IsDBNotFound(err) {
			return nil, eerrs.ErrAccountNotFound.Wrap()
		}
		return nil, err
	}
	if !pwdutil.VerifyPassword(req.Password, a.Password) {
		return nil, eerrs.ErrPassword.Wrap()
	}

	if a.EnableGoogleAuth {
		if !totp.Validate(req.Code, a.GoogleAuthSecret) {
			return nil, eerrs.ErrGoogleAuthCode.Wrap()
		}
	} else {
		if req.Code != "123456" {
			return nil, eerrs.ErrGoogleAuthCode.Wrap()
		}
	}

	adminToken, err := o.CreateToken(ctx, &admin.CreateTokenReq{UserID: a.UserID, UserType: constant.AdminUser})
	if err != nil {
		return nil, err
	}
	return &admin.LoginResp{
		AdminUserID:  a.UserID,
		AdminAccount: a.Account,
		AdminToken:   adminToken.Token,
		Nickname:     a.Nickname,
		FaceURL:      a.FaceURL,
		Level:        a.Level,
	}, nil
}

func (o *adminServer) ChangePassword(ctx context.Context, req *admin.ChangePasswordReq) (*admin.ChangePasswordResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	a, err := o.Database.GetAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := pwdutil.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}
	update, err := ToDBAdminUpdatePassword(hashedPassword)
	if err != nil {
		return nil, err
	}
	if err := o.Database.UpdateAdmin(ctx, a.UserID, update); err != nil {
		return nil, err
	}
	return &admin.ChangePasswordResp{}, nil
}

func (o *adminServer) genUserID() string {
	const l = 10
	data := make([]byte, l)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9]
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}

func (o *adminServer) CheckSuperAdmin(ctx context.Context) error {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return err
	}

	adminUser, err := o.Database.GetAdminUserID(ctx, userID)
	if err != nil {
		return err
	}

	if adminUser.Level != constant.AdvancedUserLevel {
		return errs.ErrNoPermission.Wrap()
	}
	return nil
}

func (o *adminServer) GetGoogleAuth(ctx context.Context, req *admin.GetGoogleAuthReq) (*admin.GetGoogleAuthResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}

	adminInfo, err := o.Database.GetAdminUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if adminInfo.GoogleAuthSecret == "" {
		secret, err := generateGoogleAuthSecret()
		if err != nil {
			return nil, err
		}
		adminInfo.GoogleAuthSecret = secret
		// admin => map[string]any
		update := make(map[string]any)
		update["google_auth_secret"] = secret
		if err := o.Database.UpdateAdmin(ctx, adminInfo.UserID, update); err != nil {
			return nil, err
		}
	}

	qrCodeData := generateGoogleAuthQRCode(adminInfo.Nickname, adminInfo.GoogleAuthSecret)

	return &admin.GetGoogleAuthResp{
		Secret:    adminInfo.GoogleAuthSecret,
		QrCodeUrl: qrCodeData,
	}, nil
}

func (o *adminServer) VerifyGoogleAuth(ctx context.Context, req *admin.VerifyGoogleAuthReq) (*admin.VerifyGoogleAuthResp, error) {
	userID, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}

	adminInfo, err := o.Database.GetAdminUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !totp.Validate(req.Code, adminInfo.GoogleAuthSecret) {
		return nil, eerrs.ErrGoogleAuthCode.Wrap()
	}

	adminInfo.EnableGoogleAuth = true
	update := make(map[string]any)
	update["enable_google_auth"] = true
	if err := o.Database.UpdateAdmin(ctx, adminInfo.UserID, update); err != nil {
		return nil, err
	}

	return &admin.VerifyGoogleAuthResp{}, nil
}

func generateGoogleAuthSecret() (string, error) {
	secret := make([]byte, 10)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

func generateGoogleAuthQRCode(nickName, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=https://43.154.73.22:58602",
		"飞宏IM", // 改为你的应用名称
		nickName,
		secret)
}
