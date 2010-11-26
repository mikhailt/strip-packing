package main

import (
  "cmath"
)

type Bin struct {
  Rect
  bins Bin[]
  top int
  t int
}

// Packs rectangles from slice 'a' to the strip starting from lower bound 'be'
// according to Kuzyurin-Pospelov's basic algorithm. 
// Returns an upper bound of resulting alignment. This algo in not quite 
// on-line, it uses number of all the rectangles.
type Kp1Algo struct {
  delta, u float64
  d int
}

func (v *Kp1Algo) Pack(a []Rect, be float64) float64 {
  var top float64 = be
  n := len(a)
  bins = make([]Bin, n)
  nbins = 0
  var j int
  
  v.delta = real(cmath.Pow(cmplx(float64(n), 0), -1/3))
  v.u = real(cmath.Pow(cmplx(float64(n), 0), 1/3))
  v.d = int(1 / (2 * v.delta))
  for _, r := range(a) {
    if r.w > (1 - v.delta) {
      r.y = top
      r.x = 0
      top += r.h
      continue
    }
    // Determining type of the current rectangle.
    j = 0
    for y := 1; y <= v.d; y++ {
      if r.w <= (v.delta * y) {
        j = y
        break
      }
    }
    if 0 == j {
      for y := v.d; y >= 1; y-- {
        if r.w <= (1 - v.delta * y) {
          j = v.d * 2 - y + 1
          break
        }
      }
    }
    // Finding suitable opened bin for current rectangle.
    bin_found := false
    for y := 0; y < nbins; y++ {
      if bins[y].t != j {
        continue
      }
      if (bins[y].h - bins[y].top) >= r.w {
        r.x = bin[y].x
        r.y = bin[y].y + bin[y].top
        bin[y].top += r.h
        bind_found = true
        break
      }
    }
    if bin_found {
      continue
    }
    // Opening new bin for current rectangle.
    var jj int
    if j > v.d {
      jj = 2 * v.d - j + 1
    } else {
      jj = j
    }
    
  }
  return top
}

func (v *Kp1Algo) AddBin(x, y, h, w float64) {
  v.bins[v.nbins].x = x
  v.bins[v.nbins].y = y
  v.bins[v.nbins].h = h
  v.bins[v.nbins].w = w
  v.nbins++
}