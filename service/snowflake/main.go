package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"hsm/service/snowflake/handler"
	"log"

	snowflake "hsm/service/snowflake/proto/snowflake"
)

func main() {
	// 初始化服务发现 etcd
	consulReg := etcdv3.NewRegistry(func(options *registry.Options) {
		options.Addrs=[]string{
			"192.168.5.88:2379",
			"192.168.5.138:2379",
			"192.168.3.31:2379",
		}
	})

	// 初始化服务发现 consul
	//consulReg := consul.NewRegistry(registry.Addrs("192.168.5.88:8500"))

	// 初始化micro服务对象，指定consul为服务发现
	service := micro.NewService(
		micro.Name("go.micro.srv.snowflake"),
		micro.Registry(consulReg),
		micro.Version("latest"),
		micro.Address(":50002"),
	)

	// Register Handler
	snowflake.RegisterSnowflakeHandler(service.Server(), new(handler.Snowflake))

	fmt.Println("启动服务中...")
	// Run service
	if err := service.Run(); err != nil {
		log.Panicln(err)
	}
}
