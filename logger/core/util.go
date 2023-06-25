package core

import "math"

func (r1 Reading) MaxPcntDifference(r2 Reading) float64 {
	p1 := math.Abs(r2.Temperatures[0]-r1.Temperatures[0]) / r1.Temperatures[0]
	p2 := math.Abs(r2.Temperatures[1]-r1.Temperatures[1]) / r1.Temperatures[1]
	return math.Max(p1, p2)
}
