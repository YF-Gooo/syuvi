package main

import (
	"fmt"
	. "syuvi/core"
	module "syuvi/modules"
	"time"
)

func main() {
	core := NewCore(100)
	c1 := module.NewDemoModule("c1", "1")
	c2 := module.NewDemoModule("c2", "2")
	core.RegisterModule("c1", c1)
	core.RegisterModule("c2", c2)
	if err := core.Start(); err != nil {
		fmt.Printf("start error %v\n", err)
	}
	fmt.Println(core.Start())
	time.Sleep(time.Second * 2)
	core.Stop()
	core.Destory()
}
