package main

import "fmt"

func main() {
	for i := 1; i < 50; i++ {
		defer func() {
			fmt.Println("--===========")
		}()
		fmt.Println(i)
	}
}
