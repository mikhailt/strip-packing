package main

import (
	"fmt"
)

func Validate(a [][]Rect) bool {
	for y := 0; y < len(a); y++ {
		if false == ValidateStrip(a[y]) {
			return false
		}
	}
	return true
}

// Validates that rectangles from a are packed into strip correctly. I.e. 
// without overlapping and strictly inside half-infine strip of width 1.
// Returns true/false depending on correctness. 
func ValidateStrip(a []Rect) bool {
	cnt := len(a)
	for y := 0; y < cnt; y++ {
		if !rect_inside_strip(&a[y]) {
			println("Following rectangle intersects strip borders")
			fmt.Printf("x = %0.9v, w = %0.9v, y = %0.9v, h = %0.9v\n",
				a[y].x, a[y].w, a[y].y, a[y].h)
			return false
		}
	}
	for y := 0; y < cnt-1; y++ {
		for j := y + 1; j < cnt; j++ {
			if rects_overlap(&a[y], &a[j]) {
				println("Following 2 rectangles overlap")
				fmt.Printf("x = %0.9v, w = %0.9v, y = %0.9v, h = %0.9v\n",
					a[y].x, a[y].w, a[y].y, a[y].h)
				fmt.Printf("x = %0.9v, w = %0.9v, y = %0.9v, h = %0.9v\n",
					a[j].x, a[j].w, a[j].y, a[j].h)
				return false
			}
		}
	}
	return true
}

func rect_inside_strip(r *Rect) bool {
	return !float64_less(r.x, 0) &&
		!float64_less(1, r.x+r.w) &&
		!float64_less(r.y, 0)
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
