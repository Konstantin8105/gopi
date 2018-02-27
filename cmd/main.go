package main

import (
	"fmt"
	"time"

	pi "github.com/Konstantin8105/gopi"
)

func main() {
	service := pi.NewService()
	service.Start()
	for i := 0; i < 60; i++ {
		time.Sleep(time.Second)
		fmt.Println("Result :", service.Result())
	}
	service.Stop()
	fmt.Println("Result :", service.Result())
}
