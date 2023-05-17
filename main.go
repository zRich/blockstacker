package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zRich/blockstacker/wrapper"
)

const (
	// ConfigFile = "./conf/client/sdk_config_org1_client1.yml"
	// ConfigFile = "./conf/client/sdk_config.yml"
	ConfigFile = "./conf/client/sdk_config_user1.yml"
)

func main() {

	r := gin.New()

	// enable cors

	r.Use(Cors())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/version", func(ctx *gin.Context) {
		cmClient, err := wrapper.CreateCMClientWithConfig(ConfigFile)
		if err != nil {
			log.Fatalln(err)
		}
		version, err := cmClient.GetChainMakerServerVersion()
		if err != nil {
			log.Fatalln(err)
		}
		ctx.JSON(200, gin.H{
			"version": version,
		})
	})
	r.Run()
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
