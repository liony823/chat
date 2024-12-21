package admin

import "context"

type BucketConfig struct {
	Key     string                 `bson:"key"`
	Options map[string]interface{} `bson:"options"`
}

func (BucketConfig) TableName() string {
	return "bucket_config"
}

type BucketConfigInterface interface {
	Set(ctx context.Context, options map[string]interface{}) error
	Get(ctx context.Context) (map[string]interface{}, error)
}
