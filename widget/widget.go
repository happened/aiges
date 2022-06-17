package widget

import (
	"github.com/xfyun/aiges/model"
	"github.com/xfyun/aiges/service"
)

type WidgetInner interface {
	Register(srv *service.EngService) (err error)
	Version() (version string)
}

func NewWidget() WidgetInner {
	return &GoWidget{
		eng: model.NewEngineImpl(),
	}
}
