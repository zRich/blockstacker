package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zRich/cm-api-server/src/client"
	"github.com/zRich/cm-api-server/src/common"
	"github.com/zRich/cm-api-server/src/logger"
)

const (
	DEFAULT_CONFIG = "/cm-api-server/config/sdk_config.yaml"
)

func main() {

	logger.SetLogConfig(logger.DefaultLogConfig())

	var configFile string

	// parentDir := path.Join("..")

	// // 拼接config子目录路径
	// configDir := path.Join(parentDir, "config")

	// // 读取config.yaml文件
	// configFile := path.Join(configDir, "config.yaml")

	configFile = configYml()
	// common.Log.Info("configFile: ", configFile)
	if configFile == "" {
		common.Log.Info("can not find param [--config], will use default")
		configFile = DEFAULT_CONFIG
	}

	client, err := client.CreateCMClientWithConfig(configFile)
	if err != nil {
		common.Log.Infof("creating client failed. err : %s", err.Error())
		common.Log.Infof("config : %s", err.Error())
		common.Log.Error(err.Error())
		return
	}

	r := gin.New()

	// enable cors

	r.Use(Cors())

	gin.SetMode(gin.ReleaseMode)

	r.POST("/invoke", func(ctx *gin.Context) {
		// 从请求中获取参数 body 为json格式, 类型为 InvokeContractListParams
		var body common.InvokeContractListParams

		err := ctx.ShouldBindJSON(&body)
		common.Log.Infof("Invoke request received. params : %v", body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		txRes, err := client.InvokeContract(body.ContractName, body.MethodName, "", common.ConvertToPbKeyValues(body.Parameters), -1, true)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message":     txRes.ContractResult.Message,
			"tx_id":       txRes.TxId,
			"blockHeight": txRes.TxBlockHeight,
			"extra_data":  txRes.ContractResult.GetContractEvent(),
			"raw":         txRes,
		})
	})

	r.GET("/query", func(ctx *gin.Context) {
		// 从请求中获取参数 body 为json格式, 类型为 InvokeContractListParams

		var body common.InvokeContractListParams

		err := ctx.ShouldBindJSON(&body)
		common.Log.Infof("Query request received. params : %v", body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		txRes, err := client.QueryContract(body.ContractName, body.MethodName, common.ConvertToPbKeyValues(body.Parameters), -1)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message":     txRes.ContractResult.Message,
			"tx_id":       txRes.TxId,
			"blockHeight": txRes.TxBlockHeight,
			"extra_data":  txRes.ContractResult.GetContractEvent(),
			"result":      string(txRes.ContractResult.Result),
			"raw":         txRes,
		})
	})

	// license address 0.0.0.0, port 8080
	r.Run("0.0.0.0:8080")
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

func configYml() string {
	configPath := flag.String("config", "", "sdk_config.yaml's path")
	flag.Parse()
	return *configPath
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
