package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	ling "hsm/service/ling/proto/ling"
)

type Ling struct{}

func (e *Ling) Handle(ctx context.Context, msg *ling.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *ling.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
