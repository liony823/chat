package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
)

const (
	stealthUser = "STEALTH_USER:"
)

type UserCacheInterface interface {
	SetStealthUser(ctx context.Context, userID string, stealth int64) error
	GetStealthUser(ctx context.Context, userID string) (int64, error)
	DelStealthUser(ctx context.Context, userID string) error
}

type UserCacheRedis struct {
	rdb           redis.UniversalClient
	stealthExpire time.Duration
}

func NewUserCacheRedis(rdb redis.UniversalClient) *UserCacheRedis {
	return &UserCacheRedis{rdb: rdb, stealthExpire: time.Hour * 24 * 365}
}

func (c *UserCacheRedis) SetStealthUser(ctx context.Context, userID string, stealth int64) error {
	key := stealthUser + userID
	return errs.Wrap(c.rdb.Set(ctx, key, strconv.FormatInt(stealth, 10), c.stealthExpire).Err())
}

func (c *UserCacheRedis) GetStealthUser(ctx context.Context, userID string) (int64, error) {
	key := stealthUser + userID
	stealth, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, errs.Wrap(err)
	}
	return strconv.ParseInt(stealth, 10, 64)
}

func (c *UserCacheRedis) DelStealthUser(ctx context.Context, userID string) error {
	key := stealthUser + userID
	return errs.Wrap(c.rdb.Del(ctx, key).Err())
}
