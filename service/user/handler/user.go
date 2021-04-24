package handler

import (
	"context"
	user "hsm/service/user/proto/user"
)

type User struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *User) Call(ctx context.Context, req *user.Request, rsp *user.Response) error {
	rsp.Name = "Hello " + req.Username
	return nil
}

func (e *User) Login(ctx context.Context, req *user.Request, rsp *user.Response) error {
	rsp.Name = "Hello " + req.Username
	rsp.Id = 100
	rsp.Status = 1
	rsp.Phone = "13200000000"
	rsp.Role = 1
	return nil
}
