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

package database

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/db/tx"
	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/cache"
	admindb "github.com/openimsdk/chat/pkg/common/db/model/admin"
	"github.com/openimsdk/chat/pkg/common/db/model/chat"
	"github.com/openimsdk/chat/pkg/common/db/table/admin"
	chatdb "github.com/openimsdk/chat/pkg/common/db/table/chat"
)

type ChatDatabaseInterface interface {
	GetUser(ctx context.Context, userID string) (account *chatdb.Account, err error)
	UpdateUseInfo(ctx context.Context, userID string, attribute map[string]any, updateCred, delCred []*chatdb.Credential) (err error)
	FindAttribute(ctx context.Context, userIDs []string) ([]*chatdb.Attribute, error)
	FindAttributeByAccount(ctx context.Context, accounts []string) ([]*chatdb.Attribute, error)
	TakeAttributeByPhone(ctx context.Context, areaCode string, phoneNumber string) (*chatdb.Attribute, error)
	TakeAttributeByEmail(ctx context.Context, Email string) (*chatdb.Attribute, error)
	TakeAttributeByAccount(ctx context.Context, account string) (*chatdb.Attribute, error)
	TakeAttributeByAddress(ctx context.Context, address string) (*chatdb.Attribute, error)
	TakeAttributeByUserID(ctx context.Context, userID string) (*chatdb.Attribute, error)
	TakeAccount(ctx context.Context, userID string) (*chatdb.Account, error)
	TakeCredentialByAccount(ctx context.Context, account string) (*chatdb.Credential, error)
	TakeCredentialsByUserID(ctx context.Context, userID string) ([]*chatdb.Credential, error)
	TakeLastVerifyCode(ctx context.Context, account string) (*chatdb.VerifyCode, error)
	Search(ctx context.Context, normalUser int32, keyword string, gender int32, pagination pagination.Pagination) (int64, []*chatdb.Attribute, error)
	SearchUser(ctx context.Context, keyword string, userIDs []string, genders []int32, pagination pagination.Pagination) (int64, []*chatdb.Attribute, error)
	CountVerifyCodeRange(ctx context.Context, account string, start time.Time, end time.Time) (int64, error)
	AddVerifyCode(ctx context.Context, verifyCode *chatdb.VerifyCode, fn func() error) error
	UpdateVerifyCodeIncrCount(ctx context.Context, id string) error
	DelVerifyCode(ctx context.Context, id string) error
	RegisterUser(ctx context.Context, register *chatdb.Register, account *chatdb.Account, attribute *chatdb.Attribute, credentials []*chatdb.Credential) error
	LoginRecord(ctx context.Context, record *chatdb.UserLoginRecord, verifyCodeID *string) error
	UpdatePassword(ctx context.Context, userID string, password string) error
	UpdatePasswordAndDeleteVerifyCode(ctx context.Context, userID string, password string, codeID string) error
	NewUserCountTotal(ctx context.Context, before *time.Time) (int64, error)
	UserLoginCountTotal(ctx context.Context, before *time.Time) (int64, error)
	UserLoginCountRangeEverydayTotal(ctx context.Context, start *time.Time, end *time.Time) (map[string]int64, int64, error)
	DelUserAccount(ctx context.Context, userIDs []string) error

	// OWL 新加

	// ########## User ##########
	GetAllUserID(ctx context.Context, pagination pagination.Pagination) (int64, []string, error)
	GetAccount(ctx context.Context, userID string) (*chatdb.Account, error)
	GetAttribute(ctx context.Context, userID string) (*chatdb.Attribute, error)
	GetAttributeByAccount(ctx context.Context, account string) (*chatdb.Attribute, error)

	// ########## Group ##########
	GetGroupFromContact(ctx context.Context, userID string) (*chatdb.Contact, error)
	DeleteGroupFromContact(ctx context.Context, userID string, groupIDs []string) error
	SaveGroupToContact(ctx context.Context, userID string, groupIDs []string) error

	// ########## Post ##########
	CreatePost(ctx context.Context, post []*chatdb.PostDB) error
	UpdatePost(ctx context.Context, postID string, data map[string]any) error
	DeletePost(ctx context.Context, postIDs []string) error
	GetCommentPostIDsByUser(ctx context.Context, userID string) ([]string, error)
	GetFollowedUserIDs(ctx context.Context, userID string) ([]string, error)
	GetSubscriberUserIDs(ctx context.Context, userID string) ([]string, error)
	GetPostByForwardPostID(ctx context.Context, userID, forwardPostID string) (*chatdb.Post, error)
	GetCommentPostsByPostID(ctx context.Context, cursor int64, postID string, count int64) ([]*chatdb.Post, string, error)
	GetPostsByCursorAndUserIDs(ctx context.Context, cursor int64, userIDs []string, count int64) ([]*chatdb.Post, string, error)
	GetPostsByCursorAndUser(ctx context.Context, cursor int64, userID string, count int64) ([]*chatdb.Post, string, error)
	GetPostByID(ctx context.Context, postID string) (*chatdb.Post, error)

	// ########## UserPostRelation ##########
	GetUserPostRelation(ctx context.Context, userID, postID string) (*chatdb.UserPostRelation, error)
	CreateUserPostRelation(ctx context.Context, relations []*chatdb.UserPostRelation) error
	UpdateUserPostRelation(ctx context.Context, userID, postID string, data map[string]any) error
	DeleteUserPostRelation(ctx context.Context, userID, postID string) error
	GetPostIDsByLike(ctx context.Context, userID string) ([]string, error)
	GetPostIDsByCollect(ctx context.Context, userID string) ([]string, error)
	GetPinnedPostByUserID(ctx context.Context, userID string) (*chatdb.Post, error)
	GetPostsByCursorAndPostIDs(ctx context.Context, cursor int64, postIDs []string, count int64) ([]*chatdb.Post, string, error)

	// ########## Stealth ##########
	SetStealthUser(ctx context.Context, userID string, stealth int64) error
	GetStealthUser(ctx context.Context, userID string) (int64, error)
	DelStealthUser(ctx context.Context, userID string) error
}

func NewChatDatabase(cli *mongoutil.Client, rdb redis.UniversalClient) (ChatDatabaseInterface, error) {
	register, err := chat.NewRegister(cli.GetDB())
	if err != nil {
		return nil, err
	}
	account, err := chat.NewAccount(cli.GetDB())
	if err != nil {
		return nil, err
	}
	attribute, err := chat.NewAttribute(cli.GetDB())
	if err != nil {
		return nil, err
	}
	credential, err := chat.NewCredential(cli.GetDB())
	if err != nil {
		return nil, err
	}
	userLoginRecord, err := chat.NewUserLoginRecord(cli.GetDB())
	if err != nil {
		return nil, err
	}
	verifyCode, err := chat.NewVerifyCode(cli.GetDB())
	if err != nil {
		return nil, err
	}
	forbiddenAccount, err := admindb.NewForbiddenAccount(cli.GetDB())
	if err != nil {
		return nil, err
	}

	post, err := chat.NewPost(cli.GetDB())
	if err != nil {
		return nil, err
	}

	userPostRelation, err := chat.NewUserPostRelation(cli.GetDB())
	if err != nil {
		return nil, err
	}

	contact, err := chat.NewContact(cli.GetDB())
	if err != nil {
		return nil, err
	}

	return &ChatDatabase{
		tx:               cli.GetTx(),
		register:         register,
		account:          account,
		attribute:        attribute,
		credential:       credential,
		userLoginRecord:  userLoginRecord,
		verifyCode:       verifyCode,
		forbiddenAccount: forbiddenAccount,

		post:             post,
		userPostRelation: userPostRelation,
		contact:          contact,
		userCache:        cache.NewUserCacheRedis(rdb),
	}, nil
}

type ChatDatabase struct {
	tx               tx.Tx
	register         chatdb.RegisterInterface
	account          chatdb.AccountInterface
	attribute        chatdb.AttributeInterface
	contact          chatdb.ContactInterface
	credential       chatdb.CredentialInterface
	userLoginRecord  chatdb.UserLoginRecordInterface
	verifyCode       chatdb.VerifyCodeInterface
	forbiddenAccount admin.ForbiddenAccountInterface

	post             chatdb.PostInterface
	userPostRelation chatdb.UserPostRelationInterface

	userCache cache.UserCacheInterface
}

func (o *ChatDatabase) SetStealthUser(ctx context.Context, userID string, stealth int64) error {
	return o.userCache.SetStealthUser(ctx, userID, stealth)
}

func (o *ChatDatabase) GetStealthUser(ctx context.Context, userID string) (int64, error) {
	return o.userCache.GetStealthUser(ctx, userID)
}

func (o *ChatDatabase) DelStealthUser(ctx context.Context, userID string) error {
	return o.userCache.DelStealthUser(ctx, userID)
}

func (o *ChatDatabase) GetUser(ctx context.Context, userID string) (account *chatdb.Account, err error) {
	return o.account.Take(ctx, userID)
}

func (o *ChatDatabase) UpdateUseInfo(ctx context.Context, userID string, attribute map[string]any, updateCred, delCred []*chatdb.Credential) (err error) {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = o.attribute.Update(ctx, userID, attribute); err != nil {
			return err
		}
		for _, credential := range updateCred {
			if err = o.credential.CreateOrUpdateAccount(ctx, credential); err != nil {
				return err
			}
		}
		if err = o.credential.DeleteByUserIDType(ctx, delCred...); err != nil {
			return err
		}
		return nil
	})
}

func (o *ChatDatabase) FindAttribute(ctx context.Context, userIDs []string) ([]*chatdb.Attribute, error) {
	return o.attribute.Find(ctx, userIDs)
}

func (o *ChatDatabase) FindAttributeByAccount(ctx context.Context, accounts []string) ([]*chatdb.Attribute, error) {
	return o.attribute.FindAccount(ctx, accounts)
}

func (o *ChatDatabase) TakeAttributeByPhone(ctx context.Context, areaCode string, phoneNumber string) (*chatdb.Attribute, error) {
	return o.attribute.TakePhone(ctx, areaCode, phoneNumber)
}

func (o *ChatDatabase) TakeAttributeByEmail(ctx context.Context, email string) (*chatdb.Attribute, error) {
	return o.attribute.TakeEmail(ctx, email)
}

func (o *ChatDatabase) TakeAttributeByAccount(ctx context.Context, account string) (*chatdb.Attribute, error) {
	return o.attribute.TakeAccount(ctx, account)
}

func (o *ChatDatabase) TakeAttributeByUserID(ctx context.Context, userID string) (*chatdb.Attribute, error) {
	return o.attribute.Take(ctx, userID)
}

func (o *ChatDatabase) TakeLastVerifyCode(ctx context.Context, account string) (*chatdb.VerifyCode, error) {
	return o.verifyCode.TakeLast(ctx, account)
}

func (o *ChatDatabase) TakeAccount(ctx context.Context, userID string) (*chatdb.Account, error) {
	return o.account.Take(ctx, userID)
}

func (o *ChatDatabase) TakeCredentialByAccount(ctx context.Context, account string) (*chatdb.Credential, error) {
	return o.credential.TakeAccount(ctx, account)
}

func (o *ChatDatabase) TakeCredentialsByUserID(ctx context.Context, userID string) ([]*chatdb.Credential, error) {
	return o.credential.Find(ctx, userID)
}

func (o *ChatDatabase) Search(ctx context.Context, normalUser int32, keyword string, genders int32, pagination pagination.Pagination) (total int64, attributes []*chatdb.Attribute, err error) {
	var forbiddenIDs []string
	if int(normalUser) == constant.NormalUser {
		forbiddenIDs, err = o.forbiddenAccount.FindAllIDs(ctx)
		if err != nil {
			return 0, nil, err
		}
	}
	total, totalUser, err := o.attribute.SearchNormalUser(ctx, keyword, forbiddenIDs, genders, pagination)
	if err != nil {
		return 0, nil, err
	}
	return total, totalUser, nil
}

func (o *ChatDatabase) SearchUser(ctx context.Context, keyword string, userIDs []string, genders []int32, pagination pagination.Pagination) (int64, []*chatdb.Attribute, error) {
	return o.attribute.SearchUser(ctx, keyword, userIDs, genders, pagination)
}

func (o *ChatDatabase) CountVerifyCodeRange(ctx context.Context, account string, start time.Time, end time.Time) (int64, error) {
	return o.verifyCode.RangeNum(ctx, account, start, end)
}

func (o *ChatDatabase) AddVerifyCode(ctx context.Context, verifyCode *chatdb.VerifyCode, fn func() error) error {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := o.verifyCode.Add(ctx, []*chatdb.VerifyCode{verifyCode}); err != nil {
			return err
		}
		if fn != nil {
			return fn()
		}
		return nil
	})
}

func (o *ChatDatabase) UpdateVerifyCodeIncrCount(ctx context.Context, id string) error {
	return o.verifyCode.Incr(ctx, id)
}

func (o *ChatDatabase) DelVerifyCode(ctx context.Context, id string) error {
	return o.verifyCode.Delete(ctx, id)
}

func (o *ChatDatabase) RegisterUser(ctx context.Context, register *chatdb.Register, account *chatdb.Account, attribute *chatdb.Attribute, credentials []*chatdb.Credential) error {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := o.register.Create(ctx, register); err != nil {
			return err
		}
		if err := o.account.Create(ctx, account); err != nil {
			return err
		}
		if err := o.attribute.Create(ctx, attribute); err != nil {
			return err
		}
		if err := o.credential.Create(ctx, credentials...); err != nil {
			return err
		}
		return nil
	})
}

func (o *ChatDatabase) LoginRecord(ctx context.Context, record *chatdb.UserLoginRecord, verifyCodeID *string) error {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := o.userLoginRecord.Create(ctx, record); err != nil {
			return err
		}
		if verifyCodeID != nil {
			if err := o.verifyCode.Delete(ctx, *verifyCodeID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (o *ChatDatabase) UpdatePassword(ctx context.Context, userID string, password string) error {
	return o.account.UpdatePassword(ctx, userID, password)
}

func (o *ChatDatabase) UpdatePasswordAndDeleteVerifyCode(ctx context.Context, userID string, password string, codeID string) error {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := o.account.UpdatePassword(ctx, userID, password); err != nil {
			return err
		}
		if codeID == "" {
			return nil
		}
		if err := o.verifyCode.Delete(ctx, codeID); err != nil {
			return err
		}
		return nil
	})
}

func (o *ChatDatabase) NewUserCountTotal(ctx context.Context, before *time.Time) (int64, error) {
	return o.register.CountTotal(ctx, before)
}

func (o *ChatDatabase) UserLoginCountTotal(ctx context.Context, before *time.Time) (int64, error) {
	return o.userLoginRecord.CountTotal(ctx, before)
}

func (o *ChatDatabase) UserLoginCountRangeEverydayTotal(ctx context.Context, start *time.Time, end *time.Time) (map[string]int64, int64, error) {
	return o.userLoginRecord.CountRangeEverydayTotal(ctx, start, end)
}

func (o *ChatDatabase) DelUserAccount(ctx context.Context, userIDs []string) error {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := o.register.Delete(ctx, userIDs); err != nil {
			return err
		}
		if err := o.account.Delete(ctx, userIDs); err != nil {
			return err
		}
		if err := o.attribute.Delete(ctx, userIDs); err != nil {
			return err
		}
		return nil
	})
}

/* OWL 新加 */

// ########## User ##########
func (o *ChatDatabase) GetAllUserID(ctx context.Context, pagination pagination.Pagination) (int64, []string, error) {
	return o.account.GetAllUserID(ctx, pagination)
}

func (o *ChatDatabase) GetAccount(ctx context.Context, userID string) (*chatdb.Account, error) {
	return o.account.Take(ctx, userID)
}

func (o *ChatDatabase) GetAttribute(ctx context.Context, userID string) (*chatdb.Attribute, error) {
	return o.attribute.Take(ctx, userID)
}

func (o *ChatDatabase) GetAttributeByAccount(ctx context.Context, account string) (*chatdb.Attribute, error) {
	return o.attribute.TakeAccount(ctx, account)
}

func (o *ChatDatabase) TakeAttributeByAddress(ctx context.Context, address string) (*chatdb.Attribute, error) {
	return o.attribute.TakeByAddress(ctx, address)
}

// ########## Group ##########
func (o *ChatDatabase) DeleteGroupFromContact(ctx context.Context, userID string, groupIDs []string) error {
	return o.contact.DeleteGroup(ctx, userID, groupIDs)
}

// SaveGroupToContact implements ChatDatabaseInterface.
func (o *ChatDatabase) SaveGroupToContact(ctx context.Context, userID string, groupIDs []string) error {
	return o.contact.AddGroup(ctx, userID, groupIDs)
}

// GetGroupFromContact implements ChatDatabaseInterface.
func (o *ChatDatabase) GetGroupFromContact(ctx context.Context, userID string) (*chatdb.Contact, error) {
	return o.contact.TakeGroups(ctx, userID)
}

// ########## Post ##########
func (o *ChatDatabase) DeletePost(ctx context.Context, postIDs []string) error {
	return o.post.Delete(ctx, postIDs)
}

func (o *ChatDatabase) GetFollowedUserIDs(ctx context.Context, userID string) ([]string, error) {
	return o.post.GetFollowedUserIDs(ctx, userID)
}

func (o *ChatDatabase) GetSubscriberUserIDs(ctx context.Context, userID string) ([]string, error) {
	return o.post.GetSubscriberUserIDs(ctx, userID)
}

func (o *ChatDatabase) GetPostByForwardPostID(ctx context.Context, userID, forwardPostID string) (*chatdb.Post, error) {
	return o.post.GetPostByForwardPostID(ctx, userID, forwardPostID)
}

func (o *ChatDatabase) GetCommentPostsByPostID(ctx context.Context, cursor int64, postID string, count int64) ([]*chatdb.Post, string, error) {
	return o.post.GetCommentPostsByPostID(ctx, cursor, postID, count)
}

func (o *ChatDatabase) GetPostsByCursorAndUserIDs(ctx context.Context, cursor int64, userIDs []string, count int64) ([]*chatdb.Post, string, error) {
	return o.post.GetPostsByCursorAndUserIDs(ctx, cursor, userIDs, count)
}

func (o *ChatDatabase) GetPostsByCursorAndUser(ctx context.Context, cursor int64, userID string, count int64) ([]*chatdb.Post, string, error) {
	return o.post.GetPostsByCursorAndUser(ctx, cursor, userID, count)
}

func (o *ChatDatabase) GetPostByID(ctx context.Context, postID string) (*chatdb.Post, error) {
	return o.post.Take(ctx, postID)
}

func (o *ChatDatabase) CreatePost(ctx context.Context, posts []*chatdb.PostDB) error {
	return o.post.Create(ctx, posts)
}

func (o *ChatDatabase) UpdatePost(ctx context.Context, postID string, data map[string]any) error {
	return o.post.UpdateByMap(ctx, postID, data)
}

func (o *ChatDatabase) GetCommentPostIDsByUser(ctx context.Context, userID string) ([]string, error) {
	return o.post.GetCommentPostIDsByUser(ctx, userID)
}

// ########## UserPostRelation ##########

func (o *ChatDatabase) GetUserPostRelation(ctx context.Context, userID, postID string) (*chatdb.UserPostRelation, error) {
	return o.userPostRelation.Take(ctx, userID, postID)
}

func (o *ChatDatabase) CreateUserPostRelation(ctx context.Context, relations []*chatdb.UserPostRelation) error {
	return o.userPostRelation.Create(ctx, relations)
}

func (o *ChatDatabase) UpdateUserPostRelation(ctx context.Context, userID, postID string, data map[string]any) error {
	return o.userPostRelation.UpdateByMap(ctx, userID, postID, data)
}

func (o *ChatDatabase) DeleteUserPostRelation(ctx context.Context, userID, postID string) error {
	return o.userPostRelation.Delete(ctx, userID, postID)
}

func (o *ChatDatabase) GetPostIDsByLike(ctx context.Context, userID string) ([]string, error) {
	return o.userPostRelation.GetLikedPostIDs(ctx, userID)
}

func (o *ChatDatabase) GetPostIDsByCollect(ctx context.Context, userID string) ([]string, error) {
	return o.userPostRelation.GetCollectedPostIDs(ctx, userID)
}

func (o *ChatDatabase) GetPinnedPostByUserID(ctx context.Context, userID string) (*chatdb.Post, error) {
	return o.post.GetPinnedPostByUserID(ctx, userID)
}

func (o *ChatDatabase) GetPostsByCursorAndPostIDs(ctx context.Context, cursor int64, postIDs []string, count int64) ([]*chatdb.Post, string, error) {
	return o.post.GetPostsByCursorAndPostIDs(ctx, cursor, postIDs, count)
}
