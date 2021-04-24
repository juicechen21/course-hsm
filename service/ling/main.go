package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"hsm/service/ling/handler"
	"hsm/service/ling/subscriber"

	ling "hsm/service/ling/proto/ling"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.ling"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	ling.RegisterLingHandler(service.Server(), new(handler.Ling))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.ling", service.Server(), new(subscriber.Ling))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.ling", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
