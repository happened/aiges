package widget

import (
	"github.com/xfyun/aiges/model"
	"github.com/xfyun/aiges/service"
)

type GoWidget struct {
	eng *model.EngineImpl
}

func (inst *GoWidget) Register(srv *service.EngService) (errInfo error) {
	// 注册"事件-行为"对
	errInfo = srv.Register(service.EventUsrInit, model.EngineInit)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrFini, model.EngineFini)
	if errInfo != nil {
		return
	}
	errInfo = srv.Register(service.EventUsrOnceExec, inst.eng.EngineOnceExec)
	if errInfo != nil {
		return
	}

	errInfo = srv.Register(service.EventUsrDebug, inst.eng.EngineDebug)
	if errInfo != nil {
		return
	}
	return
}

func (inst *GoWidget) Version() (ver string) {
	return inst.eng.EngineVersion()
}
