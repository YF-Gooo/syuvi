package module

import (
	"context"
	"errors"
	"fmt"
	. "syuvi/core"
	"time"
)

type DemoModule struct {
	evtReceiver EventReceiver
	coreCtx     context.Context
	stopChan    chan struct{}
	evtBuf      chan Event
	name        string
	content     string
}

func NewDemoModule(name string, content string) *DemoModule {
	return &DemoModule{
		stopChan: make(chan struct{}),
		evtBuf:   make(chan Event, 3),
		name:     name,
		content:  content,
	}
}

func (c *DemoModule) Init(evtReceiver EventReceiver) error {
	fmt.Println("initialize collector", c.name)
	c.evtReceiver = evtReceiver
	return nil
}

func (c *DemoModule) Start(coreCtx context.Context) error {
	fmt.Println("start collector", c.name)
	for {
		select {
		case <-coreCtx.Done():
			c.stopChan <- struct{}{}
			break
		case e := <-c.evtBuf:
			fmt.Println("server from", e)
		default:
			time.Sleep(time.Millisecond * 50)
			c.evtReceiver.OnEvent(Event{c.name, c.content})
		}
	}
}

func (c *DemoModule) Stop() error {
	fmt.Println("stop collector", c.name)
	select {
	case <-c.stopChan:
		return nil
	case <-time.After(time.Second * 1):
		return errors.New("failed to stop for timeout")
	}
}

func (c *DemoModule) Destory() error {
	fmt.Println(c.name, "released resources.")
	return nil
}

func (c *DemoModule) OnEvent(evt Event) error {
	fmt.Println(c.name, "event resources.")
	c.evtBuf <- evt
	return nil
}
