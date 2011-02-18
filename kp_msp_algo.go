package main

type Kp2MspAlgo struct {
}

func (v *Kp2MspAlgo) Pack(rects [][]Rect, be float64) float64 {
	c := make(chan float64)
	for y := 0; y < len(rects); y++ {
		go KpMspStripWorker(rects[y:y+1], be, c)
	}
	H_max := float64(0)
	for y := 0; y < len(rects); y++ {
		H_cur := <- c
		if H_cur > H_max {
			H_max = H_cur
		}
	}
	return H_max
}

func KpMspStripWorker(rects [][]Rect, be float64, c chan float64) {
	algo := new(Kp2Algo)
	c <- algo.Pack(rects, be)
}