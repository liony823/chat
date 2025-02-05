package admin

import (
	"context"
	"encoding/json"

	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/errs"
)

func (o *adminServer) GetBucketConfig(ctx context.Context, req *admin.GetBucketConfigReq) (*admin.GetBucketConfigResp, error) {
	conf, err := o.Database.GetBucketConfig(ctx)
	if err != nil {
		return nil, err
	}

	cm := make(map[string]string)
	for k, v := range conf {
		json, err := json.Marshal(v)
		if err != nil {
			cm[k] = ""
		}
		cm[k] = string(json)
	}
	return &admin.GetBucketConfigResp{Config: cm}, nil
}

func (o *adminServer) SetBucketConfig(ctx context.Context, req *admin.SetBucketConfigReq) (*admin.SetBucketConfigResp, error) {
	cm := make(map[string]interface{})
	var isActive bool
	for k, v := range req.Config {
		var object map[string]interface{}
		err := json.Unmarshal([]byte(v), &object)
		if enable, ok := object["enable"].(bool); ok && enable {
			if !isActive {
				isActive = true
			} else {
				return nil, errs.New("只能启用一种桶配置")
			}
		}
		if err != nil {
			return nil, err
		}
		cm[k] = object
	}
	err := o.Database.SetBucketConfig(ctx, cm)
	if err != nil {
		return nil, err
	}
	return &admin.SetBucketConfigResp{}, nil
}
