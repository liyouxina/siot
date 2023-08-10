package main

import (
	"fmt"
	"net"
	"strconv"
)

func sender(conn net.Conn) {
	words := "hello world!"
	conn.Write([]byte(words))
	buf := make([]byte, 4096)
	conn.Read(buf)
	fmt.Println(buf)
	fmt.Println("send over")

}

func main() {
	//server := "127.0.0.1:20000"
	//tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	//	os.Exit(1)
	//}
	//conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	//	os.Exit(1)
	//}
	//
	//fmt.Println("connect success")
	//sender(conn)
	i := 0
	for i = 0; i <= 11; i++ {
		fmt.Println("curl http://www.cnga.org.cn/cngaresource/brand.php?totalresult=1135&pageno=" + strconv.Itoa(i) + " >> " + strconv.Itoa(i) + ".html")
	}
}
