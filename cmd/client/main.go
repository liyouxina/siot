package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

const HEX_GET_DEVICE_INFO = "5aa506 %s 0051"

func sender(conn net.Conn) {
	words := "hello world!"
	conn.Write([]byte(words))
	buf := make([]byte, 4096)
	conn.Read(buf)
	fmt.Println(buf)
	fmt.Println("send over")

}

var agentPool map[string]*Agent

type Agent struct {
	Coon   net.Conn
	Status string
}

func main() {
	//aa := []byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x50}
	//var a byte
	//a = byte(0)
	//for _, v := range aa {
	//	a = a + v
	//}
	//aaaa, _ := hex.DecodeString(strconv.Itoa(int(a)))
	//fmt.Println(aaaa)
	//aa := "0123ab"
	//for _, a := range aa {
	//	ai := int(a)
	//	fmt.Println(ai)
	//}
	// 要转换的字节
	//byteData := []byte{0x48, 0x65, 0x6C, 0x6C, 0x0F} // "Hello" 的 ASCII 码
	//
	//// 将字节转换成十六进制字符串
	//hexStr := hex.EncodeToString(byteData)
	//aa, err := hex.DecodeString("5aa506ffffffff00504c")
	//fmt.Println(aa, err)
	//fmt.Printf("原始字节: %v\n", byteData)
	//fmt.Printf("转换后的十六进制字符串: %s\n", hexStr)
	//deviceDO, _ := entity.GetByDeviceId("00002711")
	//fmt.Println(deviceDO)
	//aa := fmt.Sprintf(HEX_GET_DEVICE_INFO, "00010001")
	//fmt.Println(aa)
	//aa := []byte{0x5a, 0xa5, 0x06, 0x00, 0x01, 0x00, 0x01, 0x00, 0x51}
	//for _, a := range aa {
	//	fmt.Print(string(a))
	//}
	//a, _ := getDeviceIdHex("00010001")
	//fmt.Println(*a)
	//fmt.Println(hex.EncodeToString([]byte{0x00, 0x00, 0x27, 0x11}))
	fmt.Println(fmt.Sprintf("%s %s", "asdasd"))
}

func getDeviceIdHex(deviceId string) (*string, error) {
	deviceIdInt, err := strconv.Atoi(deviceId)
	if err != nil {
		log.Warnf("转换设备号为十六进制字符串失败 %s", deviceId)
		return nil, err
	}
	hexStr := strconv.FormatInt(int64(deviceIdInt), 16)
	if len(hexStr) > 8 {
		log.Warnf("转换设备号为十六进制字符串 字符串过长 %s", deviceIdInt)
		return nil, errors.New("转换设备号为十六进制字符串 字符串过长")
	}
	for len(hexStr) < 8 {
		hexStr = "0" + hexStr
	}
	return &hexStr, nil
}
