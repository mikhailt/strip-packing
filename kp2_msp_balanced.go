package main

import (
	"container/heap"
)

type Kp2MspBalancedAlgo struct {
	frames []Bin
	Kp1Algo
}

func (v *Kp2MspBalancedAlgo) Pack(rects []Rect, xbe, ybe float64, m int) float64 {
	v.InitFrames(xbe, ybe, m)
	a := rects
	var b []Rect
	exit_flag := false
	s := 0
	for cnt := 2; !exit_flag; cnt *= 2 {
		if cnt >= len(a) {
			b = a
			exit_flag = true
		} else {
			b = a[:cnt]
			a = a[cnt:]
		}
		v.PackWithSize(b, xbe, ybe, m)
		s += len(b)
		v.RecalcFrames(rects[:s], m)
		// Dirty hack to make it possible to render bins, initial interface does
		// not have corresponding mechanism.
		if *prenderbins {
			if nil == bins_to_render {
				bins_to_render = make([]*Rect, 0)
			}
			for j := 0; j < v.d*2 + 1; j++ {
				for ; ; {
					if 0 == len(*v.bins[j]) {
						break
					}
					bins_to_render = append(bins_to_render, &heap.Pop(v.bins[j]).(*HeapElem).p.Rect)
				}
			}
		}
		// End of hack.
	}
	ans := float64(0)
	n := len(rects)
	for y := 0; y < n; y++ {
		if ans < rects[y].y + rects[y].h {
			ans = rects[y].y + rects[y].h
		}
	}
	return ans
}

func (v *Kp2MspBalancedAlgo) PackWithSize(rects []Rect, xbe, ybe float64, m int) {
	n := len(rects)
	v.Init(n)
	
	for i := 0; i < n; i++ {
		// Determine lowest strip.
		lowest := int(0)
		for y := 0; y < m; y++ {
			if v.frames[y].top < v.frames[lowest].top {
				lowest = y
			}
		}
		
		r := &rects[i]
		if r.w > (1 - v.delta) {
			PackToBin(&v.frames[lowest], r)
			continue
		}

		j := v.RectType(r)
		if v.PackToTopBin(r, j) {
			continue
		}
		// Opening pair of new bins and packing current rectangle into corresponging
		// one.
		v.PackToNewShelfInFrame(r, &v.frames[lowest], j)
		
	}
}

func (v *Kp2MspBalancedAlgo) RecalcFrames(rects []Rect, m int) {
	for y := 0; y < m; y++ {
		v.frames[y].top = 0
	}
	for y := 0; y < len(rects); y++ {
		r := &rects[y]
		t := int(r.x + 0.5*r.w)
		cur := r.y + r.h - v.frames[t].y
		if cur > v.frames[t].top {
			v.frames[t].top = cur
		}
	}
}

func (v *Kp2MspBalancedAlgo) InitFrames(xbe, ybe float64, m int) {
	v.frames = make([]Bin, m)
	for y := 0; y < m; y++ {
		v.frames[y].top = 0
		v.frames[y].w = 1
		v.frames[y].x = xbe + float64(y)
		v.frames[y].y = ybe
	}
}