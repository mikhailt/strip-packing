package main

import (
	"gdk"
	"gtk"
	"os"
)

var draw *gdk.GdkDrawable
var gc *gdk.GdkGC
var drawingarea *gtk.GtkDrawingArea

const W = 6

var H float64 = 5

const MAX_Y = 900
const MAX_X = 1212

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

func render_all(rects []Rect, m int) {
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
		draw_all(rects, m)
	},
		nil)

	vbox.Add(drawingarea)

	window.Add(vbox)
	window.Maximize()
	window.ShowAll()

	gtk.Main()
}

var pnonsolid *bool
var prenderbins *bool
var bins_to_render []*Rect = nil

func draw_all(rects []Rect, m int) {
	strip_width := W / float64(m)
	if strip_width > 1 {
		strip_width = 1
	}
	for y := 0; y < m; y++ {
		global := Rect{float64(y) * strip_width, 0, H, strip_width}
		global.Draw(false)
	}
	for _, r := range rects {
		r.w *= strip_width
		r.x *= strip_width
		r.Draw(!*pnonsolid)
	}
	if *prenderbins {
		gc.SetRgbFgColor(gdk.Color("red"))
		for _, r := range bins_to_render {
			r.w *= strip_width
			r.x *= strip_width
			r.Draw(false)
		}
		gc.SetRgbFgColor(gdk.Color("black"))
	}
}