package main

import (
	"fmt"
	"github.com/liyouxina/siot/entity"
)

func main() {
	result := entity.ListDeviceByCursor(0, 10)
	fmt.Println(result)
	result = entity.ListDeviceByCursor(10, 10)
	fmt.Println(result)
}
