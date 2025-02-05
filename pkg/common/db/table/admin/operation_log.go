package admin

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/pagination"
)

// OperationLog 管理员操作日志
type OperationLog struct {
	OperationID  string    `bson:"operation_id"`  // 操作ID
	AdminID      string    `bson:"admin_id"`      // 管理员ID
	AdminAccount string    `bson:"admin_account"` // 管理员账号
	AdminName    string    `bson:"admin_name"`    // 管理员名称
	Module       string    `bson:"module"`        // 操作模块
	Operation    string    `bson:"operation"`     // 操作说明
	Method       string    `bson:"method"`        // 请求方法
	Path         string    `bson:"path"`          // 请求路径
	IP           string    `bson:"ip"`            // 操作IP
	RequestData  string    `bson:"request_data"`  // 请求数据
	CreateTime   time.Time `bson:"create_time"`   // 创建时间
}

// TableName 表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

type OperationLogInterface interface {
	Create(ctx context.Context, operationLog []*OperationLog) error
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*OperationLog, error)
	Delete(ctx context.Context, ids []string) error
	Take(ctx context.Context, id string) (*OperationLog, error)
	Update(ctx context.Context, id string, data map[string]any) error
}
