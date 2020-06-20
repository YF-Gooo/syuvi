package module

import (
	"context"
	"fmt"
	"syuvi/event"
	"time"
)

type CacheModule struct {
	BasicModule
}
func NewCacheModule(name string, content string) *CacheModule {
	return &CacheModule{
		BasicModule{
			stopChan: make(chan struct{}, 1),
			evtBuf:   make(chan event.Event, 3),
			name:     name,
			content:  content},

	}
}
func (m *CacheModule) Start(coreCtx context.Context) error{
	fmt.Println("start CacheModule collector", m.name)
	m.state=Running
	for {
		select {
		case <-coreCtx.Done():
			goto Loop
		case <-m.stopChan:
			fmt.Println("server ",m.name,"停止运行" )
			goto Loop
		case e := <-m.evtBuf:
			fmt.Println("server from", e.Source)
		default:
			time.Sleep(time.Millisecond * 500)
			m.evtReceiver.OnEvent(event.Event{"c2", m.name, m.content})
		}
	}
Loop:
	fmt.Println("for循环外")
	return nil
}
