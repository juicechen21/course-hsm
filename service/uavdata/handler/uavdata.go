package handler

import (
	"context"
	"encoding/json"
	"fmt"

	uavdata "hsm/service/uavdata/proto/uavdata"
)

type Uavdata struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Uavdata) Call(ctx context.Context, req *uavdata.Request, rsp *uavdata.Response) error {
	v,err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println(string(v))
	rsp.Msg = string(v)
	return nil
}