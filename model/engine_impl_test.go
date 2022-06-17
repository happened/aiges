package model

import (
	"encoding/json"
	"fmt"
	"github.com/xfyun/aiges/instance"
	"github.com/xfyun/aiges/protocol"
	"io/ioutil"
	"os"
	"testing"
)

func must(err error){
	if err!=nil{
		panic(err)
	}
	return
}

func mockInput() (req *instance.ActMsg){

 	imgData:=instance.DataMeta{
		 Data: func() []byte{
			temp:=make([]byte,0)
			imgFile,err:=os.Open("E:\\project\\onnx-go\\examples\\tiny_yolov2\\data\\dog.jpg")
			must(err)
			temp,err=ioutil.ReadAll(imgFile)
			must(err)
			return temp
		 }(),
		 DataType: int(protocol.MetaDesc_IMAGE),
		 DataStatus: int(protocol.LoaderOutput_ONCE),
	}

	return &instance.ActMsg{
		DeliverData: []instance.DataMeta{imgData},
	}
}


func TestEngine(t *testing.T) {
	var (
		code int
		err error
		resp instance.ActMsg
	)
	inst:=NewEngineImpl()
	code,err=inst.EngineInit(map[string]string{})
	must(err)
	if code!=0{
		fmt.Println("errcode ",code)
		return
	}
	handle:="test001"
 	resp,code,err=inst.EngineOnceExec(handle,mockInput())
	must(err)
	var boxes []box
  	err=json.Unmarshal(resp.DeliverData[0].Data,&boxes)
	must(err)
 	  for _,v:=range boxes{
		  fmt.Println(v.Classes,v.Confidence)
	  }
}