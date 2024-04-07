package main

import (
	"fmt"
	"sync"
)

var lock sync.Mutex
var value = 0

func main() {
	lock = sync.Mutex{}
	go print()
	for i := 0; i < 30; i++ {
		go print()
	}

}

func print() {
	for value < 100 {
		fmt.Println(value)
		value = value + 1
	}
}
