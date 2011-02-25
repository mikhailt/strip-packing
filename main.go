package main

import (
	"rand"
	"time"
	"fmt"
	"flag"
	"cmath"
)

type Rect struct {
	x, y float64
	h, w float64
}

func GenerateRectangles(n int) []Rect {
	res := make([]Rect, n)
	for y := 0; y < n; y++ {
		res[y].h = rand.Float64()
		res[y].w = rand.Float64()
		res[y].x = rand.Float64()
		res[y].y = rand.Float64()
	}
	return res
}

func TotalArea(rects []Rect) float64 {
	var res float64 = 0
	for _, r := range rects {
		res += r.w * r.h
	}
	return res
}

type Algorithm interface {
	Pack(rects []Rect, xbe, ybe float64, m int) float64
}

// Returns uncovered area divided by N^(2/3).
func run(n int, render, validate bool, algo_name string, m int) (coefficient float64) {
	rects := GenerateRectangles(n)
	var algo Algorithm

	if "kp1" == algo_name {
		algo = new(Kp1Algo)
	} else if "kp2" == algo_name {
		algo = new(Kp2Algo)
	} else if "2d" == algo_name {
		algo = new(TdAlgo)
	} else if "kp2_msp" == algo_name {
		algo = new(Kp2MspAlgo)
	} else if "kp2_msp_b" == algo_name {
		algo = new(Kp2MspBalancedAlgo)
	} else {
		algo = new(Kp2Algo)
	}

	H = algo.Pack(rects, 0, 0, m)
	total_area := TotalArea(rects)
	fmt.Printf("Solution height = %0.9v\nTotal area = %0.9v\n", H, total_area)
	uncovered_area := H * float64(m) - total_area
	fmt.Printf("Uncovered area = %0.9v\n", uncovered_area)
	coefficient = uncovered_area / real(cmath.Pow(complex(float64(n), 0), (2.0/3)))

	if true == validate {
		if false == Validate(rects, m, H) {
			println("Validation: ERROR")
		} else {
			println("Validation: OK")
		}
	}

	if render {
		render_all(rects, m)
	}
	return
}

func main() {
	prender := flag.Bool("r", false, "Render resulting alignment of all the rectangles")
	prenderbins = flag.Bool("rb", false, "Render bins")
	pnonsolid = flag.Bool("ns", false, "Non solid rendering of rectangles")
	pn := flag.Int("n", 100, "Number of rectangles")
	pm := flag.Int("m", 1, "Number of strips")
	pvalidate := flag.Bool("v", false, "Validate resulting alignment")
	palgo := flag.String("a", "kp2", "Type of algorithm")
	ptimes := flag.Int("t", 1, "Number of tests")
	flag.Parse()
	rand.Seed(time.Nanoseconds())

	println("Number of rectangles = ", *pn)
	fmt.Printf("N^(2/3) = %0.9v\n\n", real(cmath.Pow(complex(float64(*pn), 0), (2.0/3))))

	var coef_s float64 = 0
	for y := 0; y < *ptimes; y++ {
		coef := run(*pn, *prender, *pvalidate, *palgo, *pm)
		coef_s += coef
	}
	fmt.Printf("\nAverage coefficient = %0.9v\n", coef_s/float64(*ptimes))
}
