package convert

import (
	"context"
	"time"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	adminpb "github.com/openimsdk/chat/pkg/protocol/admin"
)

func OperationLogDB2PB(ctx context.Context, log *admindb.OperationLog) (*adminpb.OperationLog, error) {
	return &adminpb.OperationLog{
		OperationID:  log.OperationID,
		AdminID:      log.AdminID,
		AdminAccount: log.AdminAccount,
		AdminName:    log.AdminName,
		Module:       log.Module,
		Operation:    log.Operation,
		Method:       log.Method,
		Path:         log.Path,
		IP:           log.IP,
		RequestData:  log.RequestData,
		CreateTime:   log.CreateTime.UnixMilli(),
	}, nil
}

func OperationLogPB2DB(ctx context.Context, log *adminpb.OperationLog) (*admindb.OperationLog, error) {
	return &admindb.OperationLog{
		OperationID:  log.OperationID,
		AdminID:      log.AdminID,
		AdminAccount: log.AdminAccount,
		AdminName:    log.AdminName,
		Module:       log.Module,
		Operation:    log.Operation,
		Method:       log.Method,
		Path:         log.Path,
		IP:           log.IP,
		RequestData:  log.RequestData,
		CreateTime:   time.UnixMilli(log.CreateTime),
	}, nil
}

func OperationLogPBs2DB(ctx context.Context, logs []*adminpb.OperationLog) (logPbs []*admindb.OperationLog, err error) {
	if len(logs) == 0 {
		return nil, nil
	}

	for _, log := range logs {
		logDb, err := OperationLogPB2DB(ctx, log)
		if err != nil {
			return nil, err
		}
		logPbs = append(logPbs, logDb)
	}
	return logPbs, nil
}

func OperationLogDBs2PBs(ctx context.Context, logs []*admindb.OperationLog) (logPbs []*adminpb.OperationLog, err error) {
	if len(logs) == 0 {
		return nil, nil
	}

	for _, log := range logs {
		logPb, err := OperationLogDB2PB(ctx, log)
		if err != nil {
			return nil, err
		}
		logPbs = append(logPbs, logPb)
	}
	return logPbs, nil
}
