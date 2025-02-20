package admin

import (
	"context"

	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewOperationLog(db *mongo.Database) (admindb.OperationLogInterface, error) {
	coll := db.Collection("operation_logs")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "operation_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &OperationLogMgo{coll: coll}, nil
}

type OperationLogMgo struct {
	coll *mongo.Collection
}

func (o *OperationLogMgo) Create(ctx context.Context, operationLog []*admindb.OperationLog) error {
	return mongoutil.InsertMany(ctx, o.coll, operationLog)
}

func (o *OperationLogMgo) Delete(ctx context.Context, ids []string) error {
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"operation_id": bson.M{"$in": ids}})
}

func (o *OperationLogMgo) Take(ctx context.Context, id string) (*admindb.OperationLog, error) {
	return mongoutil.FindOne[*admindb.OperationLog](ctx, o.coll, bson.M{"operation_id": id})
}

func (o *OperationLogMgo) Update(ctx context.Context, id string, data map[string]any) error {
	return mongoutil.UpdateOne(ctx, o.coll, bson.M{"operation_id": id}, bson.M{"$set": data}, false)
}

func (o *OperationLogMgo) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*admindb.OperationLog, error) {
	filter := bson.M{}
	if keyword != "" {
		filter["$or"] = bson.A{
			bson.M{"admin_name": bson.M{"$regex": keyword}},
			bson.M{"module": bson.M{"$regex": keyword}},
			bson.M{"operation": bson.M{"$regex": keyword}},
			bson.M{"path": bson.M{"$regex": keyword}},
			bson.M{"ip": bson.M{"$regex": keyword}},
			bson.M{"request_data": bson.M{"$regex": keyword}},
			bson.M{"method": bson.M{"$regex": keyword}},
		}
	}
	opts := options.Find().SetSort(bson.M{"create_time": -1})
	return mongoutil.FindPage[*admindb.OperationLog](ctx, o.coll, filter, pagination, opts)
}
