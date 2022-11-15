package core

import "math"

func (r1 Reading) MaxPcntDifference(r2 Reading) float64 {
	p1 := math.Abs(r2.Temp1-r1.Temp1) / r1.Temp1
	p2 := math.Abs(r2.Temp2-r1.Temp2) / r1.Temp2
	return math.Max(p1, p2)
}
