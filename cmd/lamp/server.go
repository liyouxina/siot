package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
)

type Resp struct {
	Msg string `json:"msg"`
}

func serve() {
	server := gin.Default()
	server.GET("/getAllAgentsWithDeviceId", getAllAgentsWithDeviceId)
	server.GET("/getAllAgentsWithSystemId", getAllAgentsWithSystemId)
	server.GET("/sendMsgBySystemId", sendMsgBySystemId)
	server.GET("/sendMsgByDeviceId", sendMsgByDeviceId)
	server.GET("/openLamp", sendMsgByDeviceId)
	server.GET("/closeLamp", sendMsgByDeviceId)
	server.GET("/openOrCloseLamp", lightControl)
	server.GET("/getDeviceInfo", getDeviceInfo)
	server.GET("/", index)
	_ = server.Run("0.0.0.0:8002")
}

func index(context *gin.Context) {
	contentByte, err := os.ReadFile("./lamp.html")
	if err != nil {
		log.Warnf("网站首页文件读取失败 %s", err.Error())
		context.JSON(200, err.Error())
		return
	}
	_, _ = context.Writer.Write(contentByte)
}

func getAllAgentsWithSystemId(context *gin.Context) {
	context.JSON(200, systemIdAgentPool)
	return
}

func getAllAgentsWithDeviceId(context *gin.Context) {
	context.JSON(200, deviceIdAgentPool)
	return
}

func sendMsgBySystemId(context *gin.Context) {
	systemId := context.Query("systemId")
	hexContent := context.Query("hex")
	agent := systemIdAgentPool[systemId]
	if agent == nil {
		context.JSON(200, Resp{
			Msg: "没找到这个设备",
		})
		return
	}
	resp, err := agent.sendMsg(hexContent)
	if err != nil {
		context.JSON(200, Resp{
			Msg: "接收返回数据出错 " + err.Error(),
		})
		return
	}
	context.JSON(200, Resp{
		Msg: *resp,
	})
}

func lightControl(context *gin.Context) {
	deviceId := context.Query("deviceId")
	light := context.Query("light")
	agent := deviceIdAgentPool[deviceId]
	if agent == nil {
		context.JSON(200, Resp{
			Msg: "没找到这个设备",
		})
		return
	}
	resp, err := agent.SingleLightControl(light)
	if err != nil {
		context.JSON(200, Resp{
			Msg: err.Error(),
		})
		return
	}
	context.JSON(200, resp)
}

func sendMsgByDeviceId(context *gin.Context) {
	deviceId := context.Query("deviceId")
	hexContent := context.Query("hex")
	agent := deviceIdAgentPool[deviceId]
	if agent == nil {
		context.JSON(200, Resp{
			Msg: "没找到这个设备",
		})
		return
	}
	resp, err := agent.sendMsg(hexContent)
	if err != nil {
		context.JSON(200, Resp{
			Msg: "接收返回数据出错 " + err.Error(),
		})
		return
	}
	context.JSON(200, Resp{
		Msg: *resp,
	})
}

func closeLamp(context *gin.Context) {
	deviceId := context.Query("deviceId")
	systemId := context.Query("systemId")
	agent := deviceIdAgentPool[deviceId]
	if agent == nil {
		agent = systemIdAgentPool[systemId]
	}
	if agent == nil {
		context.JSON(200, Resp{
			Msg: "没找到这个设备",
		})
	}

}

func openLamp(context *gin.Context) {
	deviceId := context.Query("deviceId")
	systemId := context.Query("systemId")
	agent := deviceIdAgentPool[deviceId]
	if agent == nil {
		agent = systemIdAgentPool[systemId]
	}
	if agent == nil {
		context.JSON(200, Resp{
			Msg: "没找到这个设备",
		})
	}
}

func getDeviceInfo(context *gin.Context) {
	deviceId := context.Query("deviceId")
	systemId := context.Query("systemId")
	agent := deviceIdAgentPool[deviceId]
	if agent == nil {
		agent = systemIdAgentPool[systemId]
	}
	if agent == nil {
		context.JSON(200, Resp{
			Msg: "没找到这个设备",
		})
		return
	}
	resp, err := agent.GetDeviceInfo()
	if err != nil {
		context.JSON(200, Resp{
			Msg: err.Error(),
		})
		return
	}
	context.JSON(200, resp)
}

func toJSONString(content interface{}) string {
	contentBytes, _ := json.Marshal(content)
	return string(contentBytes)
}
