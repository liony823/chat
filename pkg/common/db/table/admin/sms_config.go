package admin

import "context"

type SmsConfig struct {
	Key     string                 `bson:"key"`
	Options map[string]interface{} `bson:"options"`
}

func (SmsConfig) TableName() string {
	return "sms_config"
}

type SmsConfigInterface interface {
	Set(ctx context.Context, options map[string]interface{}) error
	Get(ctx context.Context) (map[string]interface{}, error)
}
