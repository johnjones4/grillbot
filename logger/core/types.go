package core

import (
	"context"
	"time"
)

type Device interface {
	Start(context.Context, chan error)
	GetCalibration() Calibration
	SetCalibration(Calibration)
}

type Reading struct {
	Received time.Time `json:"received"`
	Temp1    float64   `json:"temp1"`
	Temp2    float64   `json:"temp2"`
}

type Calibration struct {
	Temp1 float64 `json:"temp1"`
	Temp2 float64 `json:"temp2"`
}

type Metadata struct {
	Food   string    `json:"food"`
	Method string    `json:"method"`
	Start  time.Time `json:"start"`
}

type Listener func(Session, Reading)

type Session interface {
	NewReading(r Reading)
	GetReadings() ([]Reading, error)
	AddListener(l Listener)
	SetMetadata(m Metadata) error
	GetMetadata() (Metadata, error)
}

type XAxis struct {
	Min  time.Time
	Max  time.Time
	Size int
}

type YAxis struct {
	Min  float64
	Max  float64
	Size int
}
type Plot struct {
	Times    []time.Time
	Readings [2][]float64
}
