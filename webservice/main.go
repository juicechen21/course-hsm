package main

import (
	"github.com/gin-gonic/gin"
	"hsm/webservice/controller"
)

func main()  {
	r := gin.Default()
	r.POST("/login", controller.LonginHandler)
	r.POST("/push-device-data", controller.PushUavDataHandler)
	r.Run(":8080")
}
