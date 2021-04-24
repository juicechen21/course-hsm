package subscriber

import (
	"context"
	"log"

	user "hsm/service/user/proto/user"
)

type User struct{}

func (e *User) Handle(ctx context.Context, msg *user.Response) error {
	log.Panicln("Handler Received message: ", msg.Name)
	return nil
}

func Handler(ctx context.Context, msg *user.Response) error {
	log.Panicln("Handler Received message: ", msg.Name)
	return nil
}
