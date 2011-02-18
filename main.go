package main

import (
	"gdk"
	"gtk"
	"os"
	"rand"
	"time"
	"fmt"
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
		draw.DrawRectangle(gc, filled, conv_x(r.x), conv_y(r.y+r.h), conv_w(r.w), 1)
		return
	}
	draw.DrawRectangle(gc, filled, conv_x(r.x), conv_y(r.y+r.h), conv_w(r.w), conv_h(r.h))
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

func draw_all(a [][]Rect) {
	strip_width := W / float64(len(a))
	for y := 0; y < len(a); y++ {
		global := Rect{float64(y) * strip_width, 0, H, strip_width}
		global.Draw(false)
		for _, r := range a[y] {
			// Dirty hack! All rendering code have to be reconsidered.
			r.w *= strip_width
			r.x *= strip_width
			r.x += float64(y) * strip_width
			r.Draw(true)
		}
	}
}

func GenerateRectangles(n int, m int) [][]Rect {
	res := make([][]Rect, m)
	for y := 0; y < m; y++ {
		var one_more int
		if y >= (n % m) {
			one_more = 0
		} else {
			one_more = 1
		}
		res[y] = make([]Rect, (n / m) + one_more)
		for j := 0; j < len(res[y]); j++ {
			res[y][j].h = rand.Float64()
			res[y][j].w = rand.Float64()
			res[y][j].x = rand.Float64()
			res[y][j].y = rand.Float64()
		}
	}
	return res
}

func TotalArea(rects [][]Rect) float64 {
	var res float64 = 0
	for y := 0; y < len(rects); y++ {
		for _, r := range rects[y] {
			res += r.w * r.h
		}
	}
	return res
}

type Algorithm interface {
	Pack(a [][]Rect, be float64) float64
}

// Returns uncovered area divided by N^(2/3).
func run(n int, render, validate bool, algo_name string, m int) (coefficient float64) {
	rects := GenerateRectangles(n, m)
	var algo Algorithm

	if "kp1" == algo_name {
		algo = new(Kp1Algo)
	} else if "kp2" == algo_name {
		algo = new(Kp2Algo)
	} else if "2d" == algo_name {
		algo = new(TdAlgo)
	} else if "kp2_msp" == algo_name {
		algo = new(Kp2MspAlgo)
	} else {
		algo = new(Kp2Algo)
	}

	H = algo.Pack(rects, 0)
	total_area := TotalArea(rects)
	fmt.Printf("Solution height = %0.9v\nTotal area = %0.9v\n", H, total_area)
	uncovered_area := H * float64(m) - total_area
	fmt.Printf("Uncovered area = %0.9v\n", uncovered_area)
	coefficient = uncovered_area / real(cmath.Pow(cmplx(float64(n), 0), (2.0/3)))

	if true == validate {
		if false == Validate(rects) {
			println("Validation: ERROR")
		} else {
			println("Validation: OK")
		}
	}

	if false == render {
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
	return
}

func main() {
	prender := flag.Bool("r", false, "Render resulting alignment of all the rectangles")
	pn := flag.Int("n", 100, "Number of rectangles")
	pm := flag.Int("m", 1, "Number of strips")
	pvalidate := flag.Bool("v", false, "Validate resulting alignment")
	palgo := flag.String("a", "kp2", "Type of algorithm")
	ptimes := flag.Int("t", 1, "Number of tests")
	flag.Parse()
	rand.Seed(time.Nanoseconds())

	println("Number of rectangles = ", *pn)
	fmt.Printf("N^(2/3) = %0.9v\n\n", real(cmath.Pow(cmplx(float64(*pn), 0), (2.0/3))))

	var coef_s float64 = 0
	for y := 0; y < *ptimes; y++ {
		coef := run(*pn, *prender, *pvalidate, *palgo, *pm)
		coef_s += coef
	}
	fmt.Printf("\nAverage coefficient = %0.9v\n", coef_s/float64(*ptimes))
}
