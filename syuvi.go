package main

import (
	"fmt"
	. "syuvi/core"
	module "syuvi/modules"
	"time"
)

func main() {
	core := NewCore(100)
	//c1 := module.NewBasicModule("c1", "1")
	//c2 := module.NewBasicModule("c2", "2")
	c1 := module.NewCacheModule("c1", "1")
	c2 := module.NewCacheModule("c2", "2")
	core.RegisterModule("c1", c1)
	core.RegisterModule("c2", c2)
	if err := core.Start(); err != nil {
		fmt.Printf("start error %v\n", err)
	}
	time.Sleep(time.Second * 1)
	core.StopModule("c1")
	time.Sleep(time.Second * 5)
	core.StartModule("c1")
	time.Sleep(time.Second * 5)
	core.Stop()
	core.Destory()
}
