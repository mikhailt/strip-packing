package main

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
		// Determining type of the current rectangle and keep it at 'j'.
		j := int(0)
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
		PackToBin(&v.frames[lowest], &b1.Rect)
		PackToBin(b1, r)
		b2 := v.AddBin(v.ComplType(j))
		b2.y = b1.y
		b2.x = b1.x + b1.w
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