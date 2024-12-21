package admin

import (
	"context"
	"encoding/json"

	"github.com/liony823/tools/errs"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) GetSmsConfig(ctx context.Context, req *admin.GetSmsConfigReq) (*admin.GetSmsConfigResp, error) {
	conf, err := o.Database.GetSmsConfig(ctx)
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
	return &admin.GetSmsConfigResp{Config: cm}, nil
}

func (o *adminServer) SetSmsConfig(ctx context.Context, req *admin.SetSmsConfigReq) (*admin.SetSmsConfigResp, error) {
	cm := make(map[string]interface{})
	var isActive bool
	for k, v := range req.Config {
		var object map[string]interface{}
		err := json.Unmarshal([]byte(v), &object)
		if enable, ok := object["enable"].(bool); ok && enable {
			if !isActive {
				isActive = true
			} else {
				return nil, errs.New("只能开启一种短信配置")
			}
		}
		if err != nil {
			return nil, err
		}
		cm[k] = object
	}
	err := o.Database.SetSmsConfig(ctx, cm)
	if err != nil {
		return nil, err
	}
	return &admin.SetSmsConfigResp{}, nil
}
