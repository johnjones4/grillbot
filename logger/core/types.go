package core

import (
	"context"
	"time"
)

type Device interface {
	Start(context.Context, chan error)
}

type Reading struct {
	Received time.Time `json:"received"`
	Temp1    float64   `json:"temp1"`
	Temp2    float64   `json:"temp2"`
}

type Metadata struct {
	Food   string    `json:"food"`
	Method string    `json:"method"`
	Start  time.Time `json:"start"`
}

type Listener func(Session, Reading)

type Output interface {
	Listener() Listener
	Close()
	Start(context.Context, chan error)
}

type Session interface {
	NewReading(r Reading)
	GetReadings() ([]Reading, error)
	AddListener(l Listener)
	SetMetadata(m Metadata) error
	GetMetadata() (Metadata, error)
}
