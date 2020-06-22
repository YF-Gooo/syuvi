package server

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestServerBuild(t *testing.T) {
	s:=ServerBuild()
	// 设置优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.ListenAndServe()
	for {
		select {
		case e:=<-ServerInBuf:
			fmt.Println("server from", e.Source)
			if e.Content == "quit"{
				if err := s.Shutdown(ctx); err != nil {
					log.Fatal("Server Shutdown:", err)
				}
				log.Println("Server exiting")
				goto Loop
			}
			fmt.Println("server from", e.Content)
		default:
			time.Sleep(time.Millisecond * 1000*2)
		}
	}
	Loop:
		fmt.Println("for循环外")
}