package main

import (
	"gdk"
	"gtk"
	"os"
	"rand"
	"time"
	"fmt"
	"strconv"
	"flag"
	"cmath"
)

var draw *gdk.GdkDrawable
var gc *gdk.GdkGC
var drawingarea *gtk.GtkDrawingArea

const W = 6

var H float64 = 5

const MAX_Y = 900
const MAX_X = 1100

type Rect struct {
	x, y float64
	h, w float64
}

func (r *Rect) Draw(filled bool) {
	if 0 == conv_h(r.h) {
		draw.DrawRectangle(gc, filled, conv_x(r.x), conv_y(r.y+r.h),
			conv_w(r.w), 1)
		return
	}
	draw.DrawRectangle(gc, filled, conv_x(r.x), conv_y(r.y+r.h),
		conv_w(r.w), conv_h(r.h))
}

func conv_y(y float64) int {
	return int(((H - y) / H) * MAX_Y)
}

func conv_x(x float64) int {
	return int((x / W) * MAX_X)
}

func conv_h(h float64) int {
	return int((h / H) * MAX_Y)
}

func conv_w(w float64) int {
	return int((w / W) * MAX_X)
}

func draw_all(a []Rect) {
	global := Rect{0, 0, H, 1}
	global.Draw(false)
	for _, r := range a {
		r.Draw(true)
	}
}

func GenerateRectangles(n int) []Rect {
	res := make([]Rect, n)
	for y := 0; y < n; y++ {
		res[y].h = rand.Float64()
		res[y].w = rand.Float64()
		res[y].x = rand.Float64()
		res[y].y = rand.Float64()
	}
	return res
}

func TotalArea(rects []Rect) float64 {
	var res float64 = 0
	for _, r := range rects {
		res += r.w * r.h
	}
	return res
}

type Algorithm interface {
	Pack(a []Rect, be float64) float64
}

func main() {
	prender := flag.Bool("r", false, "Render resulting alignment of all the rectangles in the strip")
	pvalidate := flag.Bool("v", false, "Validate resulting alignment")
	palgo := flag.String("a", "kp1", "Type of algorithm")
	flag.Parse()
	rand.Seed(time.Nanoseconds())
	n, _ := strconv.Atoi(flag.Arg(0))
	rects := GenerateRectangles(n)
	var algo Algorithm
	if "kp1" == *palgo {
		algo = new(Kp1Algo)
	} else if "kp2" == *palgo {
		algo = new(Kp2Algo)
	}
	H = algo.Pack(rects, 0)
	total_area := TotalArea(rects)
	println("Number of rectangles = ", n)
	fmt.Printf("Solution height = %0.9v\nTotal area = %0.9v\n", H, total_area)
	fmt.Printf("Uncovered area = %0.9v\n", H - total_area)
	fmt.Printf("N^(2/3) = %0.9v\n", real(cmath.Pow(cmplx(float64(n), 0), (2.0 / 3))))
	
	if true == *pvalidate {
		if false == Validate(rects) {
			println("Validation: ERROR")
		} else {
			println("Validation: OK")
		}
	}
	
	if false == *prender {
		return
	}

	gtk.Init(&os.Args)
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetTitle("GTK DrawingArea")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	},
		nil)

	vbox := gtk.VBox(true, 0)
	vbox.SetBorderWidth(5)
	drawingarea = gtk.DrawingArea()

	var pixmap *gdk.GdkPixmap

	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		var allocation gtk.GtkAllocation
		drawingarea.GetAllocation(&allocation)
		draw = drawingarea.GetWindow().GetDrawable()
		pixmap = gdk.Pixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)
		gc = gdk.GC(pixmap.GetDrawable())
		gc.SetRgbFgColor(gdk.Color("white"))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
		gc.SetRgbFgColor(gdk.Color("black"))
		gc.SetRgbBgColor(gdk.Color("white"))
	},
		nil)

	drawingarea.Connect("expose-event", func() {
		if pixmap != nil {
			drawingarea.GetWindow().GetDrawable().DrawDrawable(gc,
				pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
		draw_all(rects)
	},
		nil)

	vbox.Add(drawingarea)

	window.Add(vbox)
	window.Maximize()
	window.ShowAll()

	gtk.Main()
}
