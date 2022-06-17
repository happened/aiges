package model

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"math"
	"sort"
)


type configuration struct {
	ConfidenceThreshold float64 `envconfig:"confidence_threshold" default:"0.30" required:"true"`
	ClassProbaThreshold float64 `envconfig:"proba_threshold" default:"0.90" required:"true"`
}

const (
	hSize, wSize  = 416, 416
	blockSize     = 32
	gridHeight    = 13
	gridWidth     = 13
	boxesPerCell  = 5
	numClasses    = 20
	envConfPrefix = "yolo"
)
var(
	classes = []string{"aeroplane", "bicycle", "bird", "boat", "bottle",
		"bus", "car", "cat", "chair", "cow",
		"diningtable", "dog", "horse", "motorbike", "person",
		"pottedplant", "sheep", "sofa", "train", "tv/monitor"}
	anchors     = []float64{1.08, 1.19, 3.42, 4.41, 6.63, 11.38, 9.42, 5.11, 16.62, 10.52}
	scaleFactor = float32(1) // The scale factor to resize the image to hSize*wSize
	config      configuration
)

type element struct {
	Prob  float64
	Class string
}

type byProba []element

func (b byProba) Len() int           { return len(b) }
func (b byProba) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byProba) Less(i, j int) bool { return b[i].Prob < b[j].Prob }

type box struct {
	R          image.Rectangle
	Gridcell   []int
	Confidence float64
	Classes    []element
}

type byConfidence []box

func (b byConfidence) Len() int           { return len(b) }
func (b byConfidence) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byConfidence) Less(i, j int) bool { return b[i].Confidence < b[j].Confidence }

func sigmoid(sum float32) float32 {
	return float32(1.0 / (1.0 + math.Exp(float64(-sum))))
}
func sigmoid64(sum float32) float64 {
	return 1.0 / (1.0 + math.Exp(float64(-sum)))
}
func exp(val float32) float64 {
	return math.Exp(float64(val))
}

func softmax(a []float32) []float64 {
	var sum float64
	output := make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		output[i] = math.Exp(float64(a[i]))
		sum += output[i]
	}
	for i := 0; i < len(output); i++ {
		output[i] = output[i] / sum
	}
	return output
}

func getOrderedElements(input []float64) []element {
	elems := make([]element, len(input))
	for i := 0; i < len(elems); i++ {
		elems[i] = element{
			Prob:  input[i],
			Class: classes[i],
		}
	}
	sort.Sort(sort.Reverse(byProba(elems)))
	return elems
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func drawRectangle(img *image.NRGBA, r image.Rectangle, label string) {
	col := color.RGBA{255, 0, 0, 255} // Red

	// HLine draws a horizontal line
	hLine := func(x1, y, x2 int) {
		for ; x1 <= x2; x1++ {
			img.Set(x1, y, col)
		}
	}

	// VLine draws a veritcal line
	vLine := func(x, y1, y2 int) {
		for ; y1 <= y2; y1++ {
			img.Set(x, y1, col)
		}
	}

	minX := int(float32(r.Min.X) * scaleFactor)
	maxX := int(float32(r.Max.X) * scaleFactor)
	minY := int(float32(r.Min.Y) * scaleFactor)
	maxY := int(float32(r.Max.Y) * scaleFactor)
	// Rect draws a rectangle utilizing HLine() and VLine()
	rect := func(r image.Rectangle) {
		hLine(minX, maxY, maxX)
		hLine(minX, maxY, maxX)
		hLine(minX, minY, maxX)
		vLine(maxX, minY, maxY)
		vLine(minX, minY, maxY)
	}
	addLabel(img, minX+5, minY+15, label)
	rect(r)
}

func addLabel(img *image.NRGBA, x, y int, label string) {
	col := color.NRGBA{0, 255, 0, 255}
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

// from https://medium.com/@jonathan_hui/real-time-object-detection-with-yolo-yolov2-28b1b93e2088
// 1- Sort the predictions by the Confidence scores.
// 2- Start from the top scores, ignore any current prediction if we find any previous predictions that have the same Class and IoU > 0.5 with the current prediction.
// 3- Repeat step 2 until all predictions are checked.
func sanitize(boxes []box) []box {
	sort.Sort(sort.Reverse(byConfidence(boxes)))

	for i := 1; i < len(boxes); i++ {
		if boxes[i].Confidence < config.ConfidenceThreshold {
			boxes = boxes[:i]
			break
		}
		if boxes[i].Classes[0].Prob < config.ClassProbaThreshold {
			boxes = boxes[:i]
			break
		}
		for j := i + 1; j < len(boxes); {
			iou := iou(boxes[i].R, boxes[j].R)
			if iou > 0.5 && boxes[i].Classes[0].Class == boxes[j].Classes[0].Class {
				boxes = append(boxes[:j], boxes[j+1:]...)
				continue
			}
			j++
		}
	}
	return boxes
}

// evaluate the intersection over union of two rectangles
func iou(r1, r2 image.Rectangle) float64 {
	// get the intesection rectangle
	intersection := image.Rect(
		max(r1.Min.X, r2.Min.X),
		max(r1.Min.Y, r2.Min.Y),
		min(r1.Max.X, r2.Max.X),
		min(r1.Max.Y, r2.Max.Y),
	)
	// compute the area of intersection rectangle
	interArea := area(intersection)
	r1Area := area(r1)
	r2Area := area(r2)
	// compute the intersection over union by taking the intersection
	// area and dividing it by the sum of prediction + ground-truth
	// areas - the interesection area
	return float64(interArea) / float64(r1Area+r2Area-interArea)
}

func area(r image.Rectangle) int {
	return max(0, r.Max.X-r.Min.X-1) * max(0, r.Max.Y-r.Min.Y-1)
}

