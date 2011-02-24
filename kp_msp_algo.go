package main

type Kp2MspAlgo struct {
}

func (v *Kp2MspAlgo) Pack(rects []Rect, xbe, ybe float64, m int) float64 {
	c := make(chan float64)
	n := len(rects)
	for y := 0; y < m; y++ {
		l := n / m
		if (n % m) > y {
			l++
		}
		go KpMspStripWorker(rects[:l], xbe + float64(y), ybe, c)
		rects = rects[l:]
	}
	H_max := float64(0)
	for y := 0; y < m; y++ {
		H_cur := <- c
		if H_cur > H_max {
			H_max = H_cur
		}
	}
	return H_max
}

func KpMspStripWorker(rects []Rect, xbe, ybe float64, c chan float64) {
	algo := new(Kp2Algo)
	c <- algo.Pack(rects, xbe, ybe, 1)
}