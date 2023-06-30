package ui

import (
	"bytes"
	"main/core"
	"time"

	chart "github.com/wcharczuk/go-chart/v2"
)

func GenerateChart(session core.Session) ([]byte, error) {
	readings, err := session.GetReadings()
	if err != nil {
		return nil, err
	}

	s := make([]chart.Series, 2)
	for i := range s {
		data := chart.TimeSeries{
			XValues: make([]time.Time, len(readings)),
			YValues: make([]float64, len(readings)),
		}

		for j := range readings {
			data.XValues[j] = readings[j].Received
			data.YValues[j] = readings[j].Temperatures[i]
		}

		s[i] = data
	}

	graph := chart.Chart{
		Width:  800,
		Height: 600,
		Series: s,
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeValueFormatterWithFormat("3:04:05PM"),
		},
	}

	buff := bytes.NewBuffer([]byte{})

	err = graph.Render(chart.PNG, buff)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
