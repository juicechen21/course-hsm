package main

import (
	"fmt"
	"github.com/go-micro/registry/consul"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"hsm/service/user/handler"
	"log"

	user "hsm/service/user/proto/user"
)

func main() {
	// 初始化服务发现 consul
	consulReg := consul.NewRegistry(registry.Addrs("192.168.5.88:8500"))

	// 初始化micro服务对象，指定consul为服务发现
	service := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Registry(consulReg),
		micro.Version("latest"),
		micro.Address(":50000"),
	)

	// Register Handler
	user.RegisterUserHandler(service.Server(), new(handler.User))

	fmt.Println("启动服务中...")

	// Run service
	if err := service.Run(); err != nil {
		log.Panicln(err)
	}
}
