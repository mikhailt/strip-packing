package main

import (
	"cmath"
)

type Bin struct {
	Rect
	top float64
	t   int
}

type BinList struct {
	m map[int]*Bin
}

func NewBinList() *BinList {
	bl := new(BinList)
	bl.m = make(map[int]*Bin)
	return bl
}

func (v *BinList) AddBin(b *Bin) {
	v.m[len(v.m)] = b
}

// Packs rectangles from slice 'a' to the strip starting from lower bound 'be'
// according to Kuzyurin-Pospelov's basic algorithm. 
// Returns an upper bound of resulting alignment. This algo in not quite 
// on-line, it uses number of all the rectangles.
type Kp1Algo struct {
	frame    Bin
	bins     map[int]*BinList
	delta, u float64
	d        int
}

func (v *Kp1Algo) Pack(rects [][]Rect, be float64) float64 {
	return v.PackSingleStrip(rects[0], be)
}

func (v *Kp1Algo) PackSingleStrip(a []Rect, be float64) float64 {
	v.frame.top = 0
	v.frame.y = be
	v.frame.x = 0
	v.frame.w = 1

	n := len(a)
	v.bins = make(map[int]*BinList)
	var j int
	var r *Rect

	v.delta = real(cmath.Pow(complex(float64(n), 0), (-1.0 / 3)))
	v.u = real(cmath.Pow(complex(float64(n), 0), (1.0 / 3)))
	v.d = int(1 / (2 * v.delta))

	// Initialization of bins map.
	for y := 0; y <= 2*v.d+1; y++ {
		v.bins[y] = NewBinList()
	}

	for i, _ := range a {
		r = &a[i]
		if r.w > (1 - v.delta) {
			PackToBin(&v.frame, r)
			continue
		}
		// Determining type of the current rectangle.
		j = 0
		for y := 1; y <= v.d; y++ {
			if r.w <= (v.delta * float64(y)) {
				j = y
				break
			}
		}
		if 0 == j {
			for y := v.d; y >= 1; y-- {
				if r.w <= (1 - v.delta*float64(y)) {
					j = v.d*2 - y + 1
					break
				}
			}
		}
		// Finding suitable opened bin for current rectangle.
		bin_found := false
		for _, b := range v.bins[j].m {
			if (b.h - b.top) >= r.h {
				PackToBin(b, r)
				bin_found = true
				break
			}
		}
		if bin_found {
			continue
		}
		// Opening pair of new bins and packing current rectangle into corresponging
		// one.
		b1 := v.AddBin(j)
		PackToBin(&v.frame, &b1.Rect)
		PackToBin(b1, r)
		b2 := v.AddBin(v.ComplType(j))
		b2.y = b1.y
		b2.x = b1.w
	}
	return v.frame.y + v.frame.top
}

func (v *Kp1Algo) AddBin(t int) *Bin {
	b := new(Bin)
	b.h = v.u
	b.w = v.WidthType(t)
	b.top = 0
	b.t = t
	v.bins[t].AddBin(b)
	return b
}

func PackToBin(bin *Bin, r *Rect) {
	r.x = bin.x
	r.y = bin.y + bin.top
	bin.top += r.h
}

func (v *Kp1Algo) ComplType(t int) int {
	return 2*v.d - t + 1
}

func (v *Kp1Algo) WidthType(t int) float64 {
	if t <= v.d {
		return v.delta * float64(t)
	}
	return 1 - float64(v.ComplType(t))*v.delta
}

// Algorithm on top of Kp1Algo. Kp2 applies Kp1 for consequtive subsequences of 
// length 2, 4, 8, 16, etc in on-line manner so that Kp2Algo actually does not 
// know  total number of rectangles.
type Kp2Algo struct{}

func (v *Kp2Algo) Pack(rects [][]Rect, be float64) float64 {
	a := rects[0]
	var kp1 Kp1Algo
	var b []Rect
	exit_flag := false
	var H float64 = 0

	for cnt := 2; !exit_flag; cnt *= 2 {
		if cnt > len(a) {
			exit_flag = true
			b = a
		} else {
			b = a[:cnt]
			a = a[cnt:]
		}
		H = kp1.PackSingleStrip(b, H)
	}
	return H
}
