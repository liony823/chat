package admin

import (
	"context"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) GetSigninConfig(ctx context.Context, req *adminpb.GetSigninConfigReq) (*adminpb.GetSigninConfigResp, error) {
	config, err := o.Database.GetSigninConfig(ctx)
	if err != nil {
		return nil, err
	}

	var result adminpb.SigninConfig

	if len(config.ContinueSignin) > 0 {
		continueSigninConfigs := make([]*adminpb.ContinueSigninConfig, 0, len(config.ContinueSignin))
		for _, v := range config.ContinueSignin {
			continueSigninConfigs = append(continueSigninConfigs, &adminpb.ContinueSigninConfig{
				Day:       v.Day,
				Increment: v.Increment,
			})
		}

		result.ContinueSignin = continueSigninConfigs
	}

	result.SigninType = config.SigninType
	result.Rule = config.Rule
	result.RandomMin = config.RandomMin
	result.RandomMax = config.RandomMax
	result.DailySignin = config.DailySignin

	return &adminpb.GetSigninConfigResp{Config: &result}, nil
}

func (o *adminServer) SetSigninConfig(ctx context.Context, req *adminpb.SetSigninConfigReq) (*adminpb.SetSigninConfigResp, error) {

	config := &admindb.SigninConfig{
		SigninType:  req.Config.SigninType,
		RandomMin:   req.Config.RandomMin,
		RandomMax:   req.Config.RandomMax,
		Rule:        req.Config.Rule,
		DailySignin: req.Config.DailySignin,
	}

	if len(req.Config.ContinueSignin) > 0 {
		continueSigninConfigs := make([]admindb.ContinueSigninConfig, 0, len(req.Config.ContinueSignin))
		for _, v := range req.Config.ContinueSignin {
			continueSigninConfigs = append(continueSigninConfigs, admindb.ContinueSigninConfig{
				Day:       v.Day,
				Increment: v.Increment,
			})
		}
		config.ContinueSignin = continueSigninConfigs
	}

	err := o.Database.SetSigninConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	return &adminpb.SetSigninConfigResp{}, nil
}
