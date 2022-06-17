package main

import (
	"flag"
	"fmt"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/env"
	"github.com/xfyun/aiges/service"
	"github.com/xfyun/aiges/widget"
	"os"
)

func main() {
	flag.Parse()
	env.Parse()
	if len(os.Args) < 2 {
		usage()
	}
	if *conf.CmdVer {
		fmt.Println(service.VERSION)
		return
	}

	var err error

	var aisrv service.EngService
	widgetInst := widget.NewWidget()

	// 框架初始化&逆初始化
	if err = aisrv.Init(env.AIGES_ENV_VERSION); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer aisrv.Fini()

	// 注册行为
	if err = widgetInst.Register(&aisrv); err != nil {
		fmt.Println(err.Error())
		return
	}

	// 框架运行
	if err = aisrv.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

func usage() {
	fmt.Printf("TODO:加载器参数说明\n") // TODO usage() 完善
	os.Exit(0)
}
