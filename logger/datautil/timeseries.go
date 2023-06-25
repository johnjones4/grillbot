package datautil

import (
	"main/core"
	"time"
)

func NormalizeTimeDistribution(readings []core.Reading, width int) [][]float64 {
	start := readings[0].Received
	end := readings[len(readings)-1].Received
	timeInterval := end.Sub(start) / time.Duration(width)
	data := make([][]float64, 2)
	data[0] = make([]float64, width)
	data[1] = make([]float64, width)
	j := 0
	for i := 0; i < width; i++ {
		curTime := start.Add(time.Duration(i) * timeInterval)
		aggregate := [2]float64{0, 0}
		count := 0.0
		for readings[j].Received.Before(curTime) || readings[j].Received.Equal(curTime) {
			aggregate[0] += readings[j].Temperatures[0]
			aggregate[1] += readings[j].Temperatures[1]
			count++
			j++
			if j == len(readings) {
				panic("we shouldn't get here")
			}
		}
		for k := 0; k < 2; k++ {
			if count > 0 {
				data[k][i] = aggregate[k] / count
			} else if i > 0 {
				data[k][i] = data[k][i-1]
			} else {
				data[k][i] = 0
			}
		}
	}
	return data
}
