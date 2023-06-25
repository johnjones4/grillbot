package core

import (
	"context"
	"time"
)

type Device interface {
	Start(context.Context, chan error)
	GetCalibration() [2]float64
	SetCalibration([2]float64)
}

type Reading struct {
	Received     time.Time  `json:"received"`
	Temperatures [2]float64 `json:"temperatures"`
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
