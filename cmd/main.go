package main

import (
	"fmt"
	"time"

	pi "github.com/Konstantin8105/gopi"
)

func main() {
	service := pi.NewService()
	service.Start()
	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Second)
			fmt.Println("Result :", service.Result())
		}
	}()
	time.Sleep(10 * time.Second)
	service.Stop()
	fmt.Println("Result at the end :", service.Result())
}
