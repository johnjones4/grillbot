package device

import (
	"context"
	"main/core"
	"math"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type simDevice struct {
	log         *logrus.Logger
	calibration [2]float64
	start       time.Time
	session     core.Session
	lock        sync.RWMutex
}

func NewSim(log *logrus.Logger, sess core.Session) core.Device {
	return &simDevice{
		log:         log,
		calibration: [2]float64{},
		session:     sess,
		lock:        sync.RWMutex{},
	}
}

func (s *simDevice) makeReading() core.Reading {
	elapsed := time.Since(s.start)
	return core.Reading{
		Received: time.Now(),
		Temperatures: [2]float64{
			math.Sin(elapsed.Seconds()/60)*100 + 100 + s.calibration[0],
			math.Cos(elapsed.Seconds()/60)*100 + 100 + s.calibration[1],
		},
	}
}

func (s *simDevice) Start(ctx context.Context, erchan chan error) {
	stop := ctx.Done()
	timer := time.NewTicker(time.Second)
	s.start = time.Now()
	r := s.makeReading()
	s.session.NewReading(r)
	for {
		select {
		case <-stop:
			return
		case <-timer.C:
			r := s.makeReading()
			s.session.NewReading(r)
		}
	}
}

func (s *simDevice) GetCalibration() [2]float64 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.calibration
}

func (s *simDevice) SetCalibration(calibration [2]float64) {
	s.lock.Lock()
	s.calibration = calibration
	s.lock.Unlock()
	s.log.Info("calibration is now ", calibration)
}
