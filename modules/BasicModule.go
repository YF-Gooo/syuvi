package module

import (
	"context"
	"fmt"
	"syuvi/event"
	"time"
)

const (
	Waiting = iota
	Running
)
type BasicModule struct {
	evtReceiver event.EventReceiver
	coreCtx     context.Context
	stopChan    chan struct{}
	evtBuf      chan event.Event
	name        string
	content     string
	state   int
}

func NewBasicModule(name string, content string) *BasicModule {
	return &BasicModule{
		stopChan: make(chan struct{}, 1),
		evtBuf:   make(chan event.Event, 3),
		name:     name,
		content:  content,
	}
}

func (m *BasicModule) Init(evtReceiver event.EventReceiver) error {
	fmt.Println("initialize collector", m.name)
	m.evtReceiver = evtReceiver
	return nil
}

func (m *BasicModule) Start(coreCtx context.Context) error {
	fmt.Println("start collector", m.name)
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

func (m *BasicModule) Stop() error {
	fmt.Println("stop collector", m.name)
	m.stopChan <- struct{}{}
	m.state=Waiting
	return nil
}

func (m *BasicModule) Destory() error {
	fmt.Println(m.name, "released resources.")
	return nil
}

//模块接受event并且塞入evtBuf
func (m *BasicModule) OnEvent(evt event.Event) error {
	m.evtBuf<-evt
	return nil
}

func (m *BasicModule) GetState() int{
	return m.state
}
