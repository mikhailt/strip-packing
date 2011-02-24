package main

import (
	"fmt"
)

func (v *Rect) PrintInfo() {
	fmt.Printf("x = %0.9v, w = %0.9v, y = %0.9v, h = %0.9v\n",
		v.x, v.w, v.y, v.h)
}

// Validates that rectangles are packed into strips correctly. I.e. without 
// overlapping and strictly inside one of m half-infinite strips of width 1. 
// Returns true/false depending on correctness. 
func Validate(rects []Rect, m int, H float64) bool {
	cnt := len(rects)
	for y := 0; y < cnt; y++ {
		if !rect_inside_strip(&rects[y], m, H) {
			println("Following rectangle intersects strip borders")
			rects[y].PrintInfo()
			return false
		}
	}
	for y := 0; y < cnt-1; y++ {
		for j := y + 1; j < cnt; j++ {
			if rects_overlap(&rects[y], &rects[j]) {
				println("Following 2 rectangles overlap")
				rects[y].PrintInfo()
				rects[j].PrintInfo()
				return false
			}
		}
	}
	return true
}

func rect_inside_strip(r *Rect, m int, H float64) bool {
	if float64_less(r.x, 0) || float64_less(float64(m), r.x+r.w) || 
			float64_less(r.y, 0) || float64_less(H, r.y + r.h) {
		return false
	}
	for y := 1; y < m; y++ {
		if float64_less(r.x, float64(y)) && float64_less(float64(y), r.x+r.w) {
			return false
		}
	}
	return true
}

// Returns true in case two rectangles have non-zero common area regards to 
// float64 comparison function (float64_less). False o/w.
func rects_overlap(a *Rect, b *Rect) bool {
	return segments_overlap(&Segment{a.x, a.x + a.w}, &Segment{b.x, b.x + b.w}) &&
		segments_overlap(&Segment{a.y, a.y + a.h}, &Segment{b.y, b.y + b.h})
}

type Segment struct {
	x, y float64
}

func segments_overlap(a *Segment, b *Segment) bool {
	return point_inside_segment(a.x, b) ||
		point_inside_segment(a.y, b) ||
		point_inside_segment(b.x, a) ||
		point_inside_segment(b.y, a)
}

func point_inside_segment(a float64, b *Segment) bool {
	return float64_less(b.x, a) && float64_less(a, b.y)
}

const D = 1e-8

func float64_less(a, b float64) bool {
	return (a + D) <= b
}

func float64_eq(a, b float64) bool {
	return !float64_less(a, b) && !float64_less(b, a)
}
