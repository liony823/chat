package admin

import "context"

type ContinueSigninConfig struct {
	Day       int32 `bson:"day"`       // 连签天数
	Increment int64 `bson:"increment"` // 追加金额
}

type SigninConfig struct {
	SigninType     string                 `bson:"signin_type"`     // 签到奖金发放模式: 随机发放、连签发放
	RandomMin      int32                  `bson:"random_min"`      // 随机发放最小金额，单位分
	RandomMax      int32                  `bson:"random_max"`      // 随机发放最大金额，单位分
	Rule           string                 `bson:"rule"`            // 签到规则说明
	DailySignin    int64                  `bson:"daily_signin"`    // 日签金额
	ContinueSignin []ContinueSigninConfig `bson:"continue_signin"` // 连签配置          // 是否启用
}

func (SigninConfig) TableName() string {
	return "signin_config"
}

const (
	SigninTypeRandom   = "random"   // 随机发放
	SigninTypeContinue = "continue" // 连签发放
)

type SigninConfigInterface interface {
	Set(ctx context.Context, config *SigninConfig) error
	Get(ctx context.Context) (*SigninConfig, error)
}
