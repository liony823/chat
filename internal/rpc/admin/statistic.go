package admin

import (
	"context"

	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func (o *adminServer) GetUserLoginRecord(ctx context.Context, req *admin.GetUserLoginRecordReq) (*admin.GetUserLoginRecordResp, error) {
	loginRecords, err := o.Database.GetUserLoginRecord(ctx)
	if err != nil {
		return nil, err
	}

	seenUsers := make(map[string]bool)
	records := make([]*admin.UserLoginRecord, 0, len(loginRecords))
	for _, record := range loginRecords {
		if seenUsers[record.UserID] {
			continue
		}
		seenUsers[record.UserID] = true
		records = append(records, &admin.UserLoginRecord{
			UserID:     record.UserID,
			IP:         record.IP,
			DeviceName: record.DeviceName,
			DeviceID:   record.DeviceID,
			Platform:   record.Platform,
			LoginTime:  record.LoginTime.UnixMilli(),
		})
	}
	return &admin.GetUserLoginRecordResp{
		Records: records,
	}, nil
}
