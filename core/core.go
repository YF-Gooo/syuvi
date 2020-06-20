package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"syuvi/event"
	"time"
)

const (
	Waiting = iota
	Running
)

var WrongStateError = errors.New("can not take the operation in the current state")

//核心
type Core struct {
	modules map[string]Module
	evtBuf  chan event.Event
	cancel  context.CancelFunc
	ctx     context.Context
	state   int
}

//初始化一个核心
func NewCore(sizeEvtBuf int) *Core {
	core := Core{
		modules: map[string]Module{},
		evtBuf:  make(chan event.Event, sizeEvtBuf),
		state:   Waiting,
	}

	return &core
}

//运行核心
func (core *Core) Start() error {
	if core.state != Waiting {
		return WrongStateError
	}
	core.state = Running
	core.ctx, core.cancel = context.WithCancel(context.Background())
	go core.EventProcessGroutine()
	return core.startModules()
}

//停止运行
func (core *Core) Stop() error {
	if core.state != Running {
		return WrongStateError
	}
	core.state = Waiting
	core.cancel()
	return core.stopModules()
}

//摧毁核心
func (core *Core) Destory() error {
	if core.state != Waiting {
		return WrongStateError
	}
	return core.destoryModules()
}

//注册一个模组
func (core *Core) RegisterModule(name string, module Module) error {
	if core.state != Waiting {
		return WrongStateError
	}
	core.modules[name] = module
	return module.Init(core)
}

//启动单个模组
func (core *Core) StartModule(name string) error {
	if core.modules[name].GetState() != Waiting {
		return WrongStateError
	}
	var err error
	var errs ModulesError
	var mutex sync.Mutex
	go func(name string, module Module, ctx context.Context) {
		defer func() {
			mutex.Unlock()
		}()
		err = module.Start(ctx)
		mutex.Lock()
		if err != nil {
			errs.ModuleErrors = append(errs.ModuleErrors,
				errors.New(name+":"+err.Error()))
		}
	}(name, core.modules[name], core.ctx)
	if len(errs.ModuleErrors) == 0 {
		return nil
	}
	return errs
}

//关闭单个模组
func (core *Core) StopModule(name string) error {
	if core.modules[name].GetState() != Running {
		return WrongStateError
	}
	var err error
	var errs ModulesError
	if err = core.modules[name].Stop(); err != nil {
		errs.ModuleErrors = append(errs.ModuleErrors,
			errors.New(name+":"+err.Error()))
	}
	if len(errs.ModuleErrors) == 0 {
		return nil
	}
	return errs
}

//销毁单个模组
func (core *Core) DestoryModule(name string) error {
	if core.modules[name].GetState() != Waiting {
		return WrongStateError
	}
	var err error
	var errs ModulesError
	if err = core.modules[name].Destory(); err != nil {
		errs.ModuleErrors = append(errs.ModuleErrors,
			errors.New(name+":"+err.Error()))
	}
	delete(core.modules,name)
	if len(errs.ModuleErrors) == 0 {
		return nil
	}
	return errs
}

//启动所有模组
func (core *Core) startModules() error {
	var err error
	var errs ModulesError
	var mutex sync.Mutex

	for name, module := range core.modules {
		go func(name string, module Module, ctx context.Context) {
			defer func() {
				mutex.Unlock()
			}()
			err = module.Start(ctx)
			mutex.Lock()
			if err != nil {
				errs.ModuleErrors = append(errs.ModuleErrors,
					errors.New(name+":"+err.Error()))
			}
		}(name, module, core.ctx)
	}
	if len(errs.ModuleErrors) == 0 {
		return nil
	}
	return errs
}

//停止所有模组
func (core *Core) stopModules() error {
	var err error
	var errs ModulesError
	for name, module := range core.modules {
		if err = module.Stop(); err != nil {
			errs.ModuleErrors = append(errs.ModuleErrors,
				errors.New(name+":"+err.Error()))
		}
	}
	if len(errs.ModuleErrors) == 0 {
		return nil
	}
	return errs
}

//销毁所有模组
func (core *Core) destoryModules() error {
	var err error
	var errs ModulesError
	for name, module := range core.modules {
		if err = module.Destory(); err != nil {
			errs.ModuleErrors = append(errs.ModuleErrors,
				errors.New(name+":"+err.Error()))
		}
		delete(core.modules,name)
	}
	if len(errs.ModuleErrors) == 0 {
		return nil
	}
	return errs
}

//核心收取消息
func (core *Core) OnEvent(evt event.Event) {
	core.evtBuf <- evt
}

//核心处理消息
func (core *Core) EventProcessGroutine() {
	for {
		select {
		case evt := <-core.evtBuf:
			if evt.Target == "core"{
				//核心处理事件
				fmt.Println("server from", evt)
			} else{
				//核心作为交换机，模块于模块间通讯
				if module, ok := core.modules[evt.Target];ok{
					module.OnEvent(evt)
				}
			}
		case <-core.ctx.Done():
			return
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
