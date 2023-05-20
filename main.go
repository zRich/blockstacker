package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
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

	r.GET("/getBlockHeight", func(ctx *gin.Context) {
		cmClient, err := wrapper.CreateCMClientWithConfig(ConfigFile)
		if err != nil {
			log.Fatalln(err)
		}
		height, err := cmClient.GetCurrentBlockHeight()
		if err != nil {
			log.Fatalln(err)
		}
		ctx.JSON(200, gin.H{
			"height": height,
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

	r.POST("/save", func(ctx *gin.Context) {
		cmClient, err := wrapper.CreateCMClientWithConfig(ConfigFile)
		if err != nil {
			log.Fatalln(err)
		}

		type SaveReq struct {
			FileName string `json:"file_name"`
		}

		var req SaveReq
		err = ctx.BindJSON(&req)
		if err != nil {
			log.Fatalln(err)
		}

		fileName := req.FileName

		params := []*common.KeyValuePair{
			{
				Key:   "time",
				Value: []byte(strconv.FormatInt(time.Now().Unix(), 10)),
			},
			{
				Key:   "file_name",
				Value: []byte(fileName),
			},
			{
				Key:   "file_hash",
				Value: []byte(uuid.GetUUID()),
			},
		}

		txRes, err := cmClient.InvokeContract("GoFactv520", "save", "", params, -1, true)
		if err != nil {
			log.Fatalln(err)
		}

		err = wrapper.CheckProposalRequestResp(txRes, true)

		if err != nil {
			log.Fatalln(err)
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
		cmClient, err := wrapper.CreateCMClientWithConfig(ConfigFile)
		if err != nil {
			log.Fatalln(err)
		}

		type QueryReq struct {
			FileHash string `json:"file_hash"`
		}

		var req QueryReq
		err = ctx.BindJSON(&req)
		if err != nil {
			log.Fatalln(err)
		}

		fileHash := req.FileHash

		params := []*common.KeyValuePair{
			{
				Key:   "file_hash",
				Value: []byte(fileHash),
			},
		}

		txRes, err := cmClient.QueryContract("GoFactv520", "findByFileHash", params, -1)
		if err != nil {
			log.Fatalln(err)
		}

		err = wrapper.CheckProposalRequestResp(txRes, true)

		if err != nil {
			log.Fatalln(err)
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
