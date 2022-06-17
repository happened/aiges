package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	"github.com/xfyun/aiges/instance"
	"github.com/xfyun/aiges/protocol"
	"gorgonia.org/tensor"
	"gorgonia.org/tensor/native"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
	"sync"
)

var once sync.Once
var EngineImplInst *EngineImpl

type EngineImpl struct {
	modelPath string  //模型所在路径
	graph *gorgonnx.Graph
	model    *onnx.Model
}

func NewEngineImpl() *EngineImpl {
	once.Do(func() {
		EngineImplInst=new(EngineImpl)
	})
	return EngineImplInst
}

func EngineInit(cfg map[string]string) (errNum int, errInfo error) {
	return EngineImplInst.EngineInit(cfg)
}
func EngineFini() (errNum int, errInfo error) {
	return EngineImplInst.EngineFini()
}

func (el *EngineImpl) EngineInit(cfg map[string]string) (errNum int, errInfo error) {
	fmt.Println("engine init enter ")
 	if val, ok := cfg[ModelPath]; ok {
		el.modelPath = val
	} else {
		el.modelPath = DefaultModelPath
	}
	if _, err := os.Stat(el.modelPath); err != nil && os.IsNotExist(err) {
		panic(fmt.Sprintf("%v does not exist", el.modelPath))
	}

	// Create a backend receiver
	el.graph =gorgonnx.NewGraph()
	// Create a model and set the execution backend
	el.model= onnx.NewModel(el.graph )

	// read the onnx model
	b, err := ioutil.ReadFile(el.modelPath)
	if err != nil {
		panic("read model file failed . "+err.Error())
	}
	if el.model.UnmarshalBinary(b)!=nil{
		panic("unmarshal model file failed ."+ err.Error())
	}
  	return
}
func (el *EngineImpl) EngineFini() (errNum int, errInfo error) {
	return
}
func (el *EngineImpl) EngineVersion() (ver string) {
	return "1.0.0"
}
func (el *EngineImpl) EngineOnceExec(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	//图片转input
	img,err:=jpeg.Decode(bytes.NewReader(req.DeliverData[0].Data))
	if err!=nil{
		errNum=-1
		errInfo=err
		return
	}

	// find the resize scale
	imgRescaled := image.NewNRGBA(image.Rect(0, 0, wSize, hSize))
	color := color.RGBA{0, 0, 0, 255}

	draw.Draw(imgRescaled, imgRescaled.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	var m image.Image
	if (img.Bounds().Max.X - img.Bounds().Min.X) > (img.Bounds().Max.Y - img.Bounds().Min.Y) {
		scaleFactor = float32(img.Bounds().Max.Y-img.Bounds().Min.Y) / float32(hSize)
		m = resize.Resize(0, hSize, img, resize.Lanczos3)
	} else {
		scaleFactor = float32(img.Bounds().Max.X-img.Bounds().Min.X) / float32(wSize)
		m = resize.Resize(wSize, 0, img, resize.Lanczos3)
	}
	switch m.(type) {
	case *image.NRGBA:
		draw.Draw(imgRescaled, imgRescaled.Bounds(), m.(*image.NRGBA), image.ZP, draw.Src)
	case *image.YCbCr:
		draw.Draw(imgRescaled, imgRescaled.Bounds(), m.(*image.YCbCr), image.ZP, draw.Src)
	default:
		errNum=-2
		errInfo=errors.New("unhandled type")
		return
	}
	inputT := tensor.New(tensor.WithShape(1, 3, hSize, wSize), tensor.Of(tensor.Float32))
	//err = images.ImageToBCHW(img, inputT)
	err =ImageToBCHW(imgRescaled, inputT)
	if err != nil {
		errNum=-3
		errInfo=err
		return
	}
	el.model.SetInput(0,inputT)
	if err=el.graph.Run();err!=nil{
		errNum=-4
		errInfo=err
		return
	}
	//output转图片
	tensors,oerr:=el.model.GetOutputTensors()
	if oerr!=nil{
		errNum=-5
		errInfo=oerr
		return
	}
	var rlt []box
	rlt,err=processOutput(tensors)
	if err!=nil{
		errNum=-6
		errInfo=err
		return
	}
	var b []byte
 	b,err=json.Marshal(rlt)
	 if err!=nil{
		 errNum=-7
		 errInfo=err
		 return
	 }
	meta:=instance.DataMeta{
		DataId: "boxs",
		Data: b,
		DataStatus: int(protocol.LoaderOutput_ONCE),
		DataType: int(protocol.MetaDesc_TEXT),
	}
	resp.DeliverData=make([]instance.DataMeta,0)
	resp.DeliverData=append(resp.DeliverData,meta)
 	return
}
func (el *EngineImpl) EngineDebug(handle string, req *instance.ActMsg) (resp instance.ActMsg, errNum int, errInfo error) {
	return
}
func (el *EngineImpl) EngineError(errNum int) (errInfo string) {
	return
}

func processOutput(t []tensor.Tensor) (rlt []box,err error){
	dense := t[0].(*tensor.Dense)
	if err=dense.Reshape(125, gridHeight, gridWidth);err!=nil{
	return
	}
	data, err := native.Tensor3F32(dense)
	if err != nil {
		return
	}

	var boxes = make([]box, gridHeight*gridWidth*boxesPerCell)
	var counter int
	// https://github.com/pjreddie/darknet/blob/61c9d02ec461e30d55762ec7669d6a1d3c356fb2/src/yolo_layer.c#L159
	for cx := 0; cx < gridWidth; cx++ {
		for cy := 0; cy < gridHeight; cy++ {
			for b := 0; b < boxesPerCell; b++ {
				channel := b * (numClasses + 5)
				tx := data[channel][cx][cy]
				ty := data[channel+1][cx][cy]
				tw := data[channel+2][cx][cy]
				th := data[channel+3][cx][cy]
				tc := data[channel+4][cx][cy]
				tclasses := make([]float32, 20)
				for i := 0; i < 20; i++ {
					tclasses[i] = data[channel+5+i][cx][cy]
				}
				// The predicted tx and ty coordinates are relative to the location
				// of the grid cell; we use the logistic sigmoid to constrain these
				// coordinates to the range 0 - 1. Then we add the cell coordinates
				// (0-12) and multiply by the number of pixels per grid cell (32).
				// Now x and y represent center of the bounding box in the original
				// 416x416 image space.
				// https://github.com/hollance/Forge/blob/04109c856237faec87deecb55126d4a20fa4f59b/Examples/YOLO/YOLO/YOLO.swift#L154
				x := int((float32(cx) + sigmoid(tx)) * blockSize)
				y := int((float32(cy) + sigmoid(ty)) * blockSize)
				// The size of the bounding box, tw and th, is predicted relative to
				// the size of an "anchor" box. Here we also transform the width and
				// height into the original 416x416 image space.
				w := int(exp(tw) * anchors[2*b] * blockSize)
				h := int(exp(th) * anchors[2*b+1] * blockSize)

				boxes[counter] = box{
					Gridcell:   []int{cx, cy},
					R:          image.Rect(max(y-w/2, 0), max(x-h/2, 0), min(y+w/2, wSize), min(x+h/2, hSize)),
					Confidence: sigmoid64(tc),
					Classes:    getOrderedElements(softmax(tclasses)),
				}
				counter++
			}
		}
	}
	boxes = sanitize(boxes)
	return boxes,nil
}