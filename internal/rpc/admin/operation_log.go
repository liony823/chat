package admin

import (
	"context"
	"time"

	"github.com/openimsdk/chat/pkg/common/convert"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
)

func (s *adminServer) GetOperationLog(ctx context.Context, req *adminpb.GetOperationLogReq) (*adminpb.GetOperationLogResp, error) {
	_, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	logs, err := s.Database.TakeOperationLog(ctx, req.OperationID)
	if err != nil {
		return nil, err
	}

	logPb, err := convert.OperationLogDB2PB(ctx, logs)
	if err != nil {
		return nil, err
	}

	return &adminpb.GetOperationLogResp{
		OperationID:  logPb.OperationID,
		AdminID:      logPb.AdminID,
		AdminAccount: logPb.AdminAccount,
		AdminName:    logPb.AdminName,
		Module:       logPb.Module,
		Operation:    logPb.Operation,
		Method:       logPb.Method,
		Path:         logPb.Path,
		IP:           logPb.IP,
		RequestData:  logPb.RequestData,
		CreateTime:   logPb.CreateTime,
	}, nil
}

func (s *adminServer) DeleteOperationLog(ctx context.Context, req *adminpb.DeleteOperationLogReq) (*adminpb.DeleteOperationLogResp, error) {
	_, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.Database.DeleteOperationLog(ctx, req.OperationIDs)
	if err != nil {
		return nil, err
	}

	return &adminpb.DeleteOperationLogResp{}, nil
}

func (s *adminServer) CreateOperationLog(ctx context.Context, req *adminpb.CreateOperationLogReq) (*adminpb.CreateOperationLogResp, error) {
	_, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}

	create := &admindb.OperationLog{
		OperationID:  req.OperationID,
		AdminID:      req.AdminID,
		AdminAccount: req.AdminAccount,
		AdminName:    req.AdminName,
		Module:       req.Module,
		Operation:    req.Operation,
		Method:       req.Method,
		Path:         req.Path,
		IP:           req.IP,
		RequestData:  req.RequestData,
		CreateTime:   time.UnixMilli(req.CreateTime),
	}

	err = s.Database.CreateOperationLog(ctx, []*admindb.OperationLog{create})
	if err != nil {
		return nil, err
	}

	return &adminpb.CreateOperationLogResp{}, nil
}

func (s *adminServer) SearchOperationLog(ctx context.Context, req *adminpb.SearchOperationLogReq) (*adminpb.SearchOperationLogResp, error) {
	_, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}

	total, logs, err := s.Database.SearchOperationLog(ctx, req.Keyword, req.Pagination)
	if err != nil {
		return nil, err
	}

	logPb, err := convert.OperationLogDBs2PBs(ctx, logs)
	if err != nil {
		return nil, err
	}

	return &adminpb.SearchOperationLogResp{
		Logs:  logPb,
		Total: total,
	}, nil
}
