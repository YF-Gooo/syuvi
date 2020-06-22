package module

import (
	"context"
	"fmt"
	"log"
	"syuvi/event"
	server "syuvi/modules/server"
	"time"
)

type ServerModule struct {
	BasicModule
}

func NewServerModule(name string, content string) *ServerModule {
	return &ServerModule{
		BasicModule{
			stopChan: make(chan struct{}, 1),
			evtBuf:   make(chan event.Event, 3),
			name:     name,
			content:  content},

	}
}

func (m *ServerModule) Start(coreCtx context.Context) error{
	fmt.Println("start CacheModule collector", m.name)
	m.state=Running
	s:=server.ServerBuild()
	// 设置优雅退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go s.ListenAndServe()
	for {
		select {
		case <-coreCtx.Done():
			goto Loop
		case <-m.stopChan:
			fmt.Println("server ",m.name,"停止运行" )
			goto Loop
		case e := <-m.evtBuf:
			fmt.Println("server from", e.Source)
		case e:=<-server.ServerInBuf:
			fmt.Println("server from", e.Source)
			if e.Content == "quit"{
				if err := s.Shutdown(ctx); err != nil {
					log.Fatal("Server Shutdown:", err)
				}
				log.Println("Server exiting")
				goto Loop
			}
		default:
			time.Sleep(time.Millisecond * 1000*60)
			goto Loop
		}
	}
	Loop:
		fmt.Println("for循环外")
		s.Shutdown(ctx)
		return nil
}
