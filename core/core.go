package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const (
	Waiting = iota
	Running
)

var WrongStateError = errors.New("can not take the operation in the current state")

//核心
type Core struct {
	modules map[string]Module
	evtBuf  chan Event
	cancel  context.CancelFunc
	ctx     context.Context
	state   int
}

//初始化一个核心
func NewCore(sizeEvtBuf int) *Core {
	core := Core{
		modules: map[string]Module{},
		evtBuf:  make(chan Event, sizeEvtBuf),
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
	}
	if len(errs.ModuleErrors) == 0 {
		return nil
	}
	return errs
}

//核心收取消息
func (core *Core) OnEvent(evt Event) {
	core.evtBuf <- evt
}

//核心处理消息
func (core *Core) EventProcessGroutine() {
	var evtSeg [10]Event
	for {
		for i := 0; i < 10; i++ {
			select {
			case evtSeg[i] = <-core.evtBuf:
			case <-core.ctx.Done():
				return
			}
		}
		fmt.Println(evtSeg)
	}

}
