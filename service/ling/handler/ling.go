package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	ling "hsm/service/ling/proto/ling"
)

type Ling struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Ling) Call(ctx context.Context, req *ling.Request, rsp *ling.Response) error {
	log.Log("Received Ling.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Ling) Stream(ctx context.Context, req *ling.StreamingRequest, stream ling.Ling_StreamStream) error {
	log.Logf("Received Ling.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&ling.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Ling) PingPong(ctx context.Context, stream ling.Ling_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&ling.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
