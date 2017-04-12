package main

import (
	"fmt"
	"service/base"
	"service/config"
)

func main() {
	if e := config.Load("server.config"); e != nil {
		fmt.Println(e)
		return
	}
	base.CreateNode().Start()
}
