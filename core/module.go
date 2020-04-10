package core

import (
	"context"
	"strings"
)

//扩展模组
type Module interface {
	Init(evtReceiver EventReceiver) error
	Start(coreCtx context.Context) error
	Stop() error
	Destory() error
	OnEvent(evt Event) error
}

//模组错误
type ModulesError struct {
	ModuleErrors []error
}

//模组错误输出
func (me ModulesError) Error() string {
	var strs []string
	for _, err := range me.ModuleErrors {
		strs = append(strs, err.Error())
	}
	return strings.Join(strs, ";")
}
