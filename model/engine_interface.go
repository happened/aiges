package model

import "github.com/xfyun/aiges/instance"

type Engine interface {
	EngineInit(cfg map[string]string) (errNum int, errInfo error)
	EngineFini() (errNum int, errInfo error)
	EngineVersion() (ver string)
	EngineOnceExec(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error)
	EngineDebug(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error)
	EngineError(errNum int) (errInfo string)
}
