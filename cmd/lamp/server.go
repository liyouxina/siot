package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type Resp struct {
	Msg string `json:"msg"`
}

func serve() {
	server := gin.Default()
	server.GET("/getAllAgents", getAllAgents)
	server.GET("/getAllAgentsWithDeviceId", getAllAgentsWithDeviceId)
	server.GET("/sendMsgBySystemId", sendMsgBySystemId)
	server.GET("/sendMsgByDeviceId", sendMsgByDeviceId)
	server.GET("/openLamp", sendMsgByDeviceId)
	server.GET("/closeLamp", sendMsgByDeviceId)
	server.GET("/getDeviceInfo", getDeviceInfo)
	_ = server.Run("0.0.0.0:8002")
}

func getAllAgents(context *gin.Context) {
	context.JSON(200, Resp{
		Msg: toJSONString(systemIdAgentPool),
	})
	return
}

func getAllAgentsWithDeviceId(context *gin.Context) {
	context.JSON(200, Resp{
		Msg: toJSONString(deviceIdAgentPool),
	})
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
	}

}

func toJSONString(content interface{}) string {
	contentBytes, _ := json.Marshal(content)
	return string(contentBytes)
}
