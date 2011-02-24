package main

import (
	"cmath"
)

type TdAlgo struct {
	frame Bin
	delta float64
	rects []Rect
	nrects int
}

func (v *TdAlgo) Pack(rects []Rect, xbe, ybe float64, m int) float64 {
	v.frame.top = 0
	v.frame.y = ybe
	v.frame.x = xbe
	v.frame.w = 1
	v.nrects = 0
	n := len(rects)
	v.delta = real(cmath.Pow(complex(float64(n), 0), (-1.0 / 2)))
	v.rects = make([]Rect, n)
	
	for y := 0; y < n; y++ {
		r := &rects[y]
		var best_r *Rect = nil
		var best_s float64 = 0
		var best_vertical bool
		for j := 0; j < v.nrects; j++ {
			packable, vertical := v.Packable(&v.rects[j], r)
			if packable {
				if (nil != best_r) && (best_s <= v.rects[j].Area()) {
					continue
				}
				best_r = &v.rects[j]
				best_s = best_r.Area()
				best_vertical = vertical
			}
		}
		if nil == best_r {
			for j := 0; j < v.nrects; j++ {
				packable, vertical := v.SimplePackable(&v.rects[j], r)
				if packable {
					//println("splitting")
					v.SplittingPackToRect(&v.rects[j], r, vertical)
					continue
				}
			}
			//println("on top")
			v.PackRectOnTop(r)
		} else {
			//println("sane")
			PackToRect(best_r, r, best_vertical)
		}
	}
	println(v.nrects)
	return v.frame.top + v.frame.y
}

func (v *TdAlgo) SplittingPackToRect(outer *Rect, inner *Rect, vertical bool) {
	inner.x = outer.x
	inner.y = outer.y
	side_r := v.AddRect()
	if vertical {
		side_r.x = outer.x
		side_r.y = outer.y + inner.h
		side_r.w = inner.w
		side_r.h = outer.h - inner.h
		outer.w -= inner.w
		outer.x += inner.w
	} else {
		// Horizontal.
		side_r.x = outer.x + inner.w
		side_r.y = outer.y
		side_r.w = outer.w - inner.w
		side_r.h = inner.h
		outer.h -= inner.h
		outer.y += inner.h
	}
}

func (v *TdAlgo) AddRect() *Rect {
	v.nrects++
	return &v.rects[v.nrects - 1]
}

func (v *TdAlgo) PackRectOnTop(r *Rect) {
	outer := v.AddRect()
	outer.w = 1
	outer.h = r.h
	PackToBin(&v.frame, outer)
	PackToRect(outer, r, false)
}

// Returns whether inner fits for packing into outer and corresponding vertical 
// flag. Note that inner is packable when at least one of its dimensions differ 
// from corresponding dimension of outer for no more than delta.
func (v *TdAlgo) Packable(outer *Rect, inner *Rect) (packable bool, vertical bool) {
	if ((inner.w + v.delta) >= outer.w) && (inner.w <= outer.w) && (inner.h <= outer.h) {
		return true, true
	}
	if ((inner.h + v.delta) >= outer.h) && (inner.w <= outer.w) && (inner.h <= outer.h) {
		return true, false
	}
	return false, false
}

func PackToRect(outer *Rect, inner *Rect, vertical bool) {
	if vertical {
		inner.x = outer.x
		inner.y = outer.y
		outer.y += inner.h
		outer.h -= inner.h
	} else {
		// Horizontal then.
		inner.x = outer.x
		inner.y = outer.y
		outer.x += inner.w
		outer.w -= inner.w
	}
}

func (v *Rect) Area() float64 {
	return v.x * v.y
}

func (v *TdAlgo) SimplePackable(outer *Rect, inner *Rect) (packable bool, vertical bool) {
	if (inner.w > outer.w) || (inner.h > outer.h) {
		return false, false
	}
	if (outer.h - inner.h) > (outer.w - inner.w) {
		return true, false
	}
	return true, true
}