package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wxTips/utils"
)

func webStarter() {
	fmt.Print("开启了web接口")
	//1.创建路由
	r := gin.Default()
	r.Use(utils.CorsHandler())
	//2.绑定路由规则，执行的函数
	r.GET("/qrUrl", func(c *gin.Context) {
		data := gin.H{"qrUrl": utils.QrUrl}
		c.JSON(http.StatusOK, data)
	})
	//3.监听端口，默认8080
	r.Run(":8081")
}
