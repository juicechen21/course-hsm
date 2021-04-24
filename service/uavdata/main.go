package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"hsm/service/uavdata/handler"
	"log"

	uavdata "hsm/service/uavdata/proto/uavdata"
	_ "hsm/service/uavdata/model"
)

func main() {
	// 初始化服务发现 consul
	consulReg := consul.NewRegistry(registry.Addrs("192.168.5.88:8500"))

	// 初始化micro服务对象，指定consul为服务发现
	service := micro.NewService(
		micro.Name("go.micro.srv.uavdata"),
		micro.Registry(consulReg),
		micro.Version("latest"),
		micro.Address(":50001"),
	)

	// Register Handler
	uavdata.RegisterUavdataHandler(service.Server(), new(handler.Uavdata))

	fmt.Println("启动服务中...")

	// Run service
	if err := service.Run(); err != nil {
		log.Panicln(err)
	}
}
