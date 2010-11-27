package main 

import (
  "gdk"
  "gtk"
  "os"
  "rand"
  "time"
  "fmt"
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
  draw.DrawRectangle(gc, filled, conv_x(r.x), conv_y(r.y + r.h), 
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
  for _, r := range(a) {
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
  for _, r := range(rects) {
    res += r.w * r.h
  }
  return res
}

func main() {
  rand.Seed(time.Nanoseconds())
  rects := GenerateRectangles(1000000)
  var algo Kp1Algo
  H = algo.Pack(rects, 0)
  fmt.Printf("H = %0.9v\nS = %0.9v", H, TotalArea(rects))
  return

	gtk.Init(&os.Args);
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL);
	window.SetTitle("GTK DrawingArea");
	window.Connect("destroy", func() {
		gtk.MainQuit();
	}, nil);

	vbox := gtk.VBox(true, 0);
	vbox.SetBorderWidth(5);
	drawingarea = gtk.DrawingArea();
  
  var pixmap *gdk.GdkPixmap

	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		var allocation gtk.GtkAllocation
		drawingarea.GetAllocation(&allocation);
		draw = drawingarea.GetWindow().GetDrawable()
		pixmap = gdk.Pixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24);
		gc = gdk.GC(pixmap.GetDrawable())
		gc.SetRgbFgColor(gdk.Color("white"));
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1);
		gc.SetRgbFgColor(gdk.Color("black"));
		gc.SetRgbBgColor(gdk.Color("white"));
	}, nil)

	drawingarea.Connect("expose-event", func() {
		if pixmap != nil {
			drawingarea.GetWindow().GetDrawable().DrawDrawable(gc, 
			  pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
		draw_all(rects)
	}, nil);

	vbox.Add(drawingarea);

	window.Add(vbox);
	window.Maximize();
	window.ShowAll();

	gtk.Main();
}