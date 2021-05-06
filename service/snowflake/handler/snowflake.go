package handler

import (
	"context"

	snowflake "hsm/service/snowflake/proto/snowflake"
)

type Snowflake struct{}

// GenerateOnlyId 通过雪花算法生成唯一ID  is a single request handler called via client.Call or the generated client code
func (e *Snowflake) GenerateOnlyId(ctx context.Context, req *snowflake.Request, rsp *snowflake.Response) error {
	node,err:=NewWorker(req.WorkerId)
	if err != nil {
		return err
	}
	rsp.CodeId = node.GetId()
	return nil
}
