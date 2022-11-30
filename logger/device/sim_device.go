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
	calibration core.Calibration
	start       time.Time
	session     core.Session
	lock        sync.RWMutex
}

func NewSim(log *logrus.Logger, sess core.Session) core.Device {
	return &simDevice{
		log:         log,
		calibration: core.Calibration{},
		session:     sess,
		lock:        sync.RWMutex{},
	}
}

func (s *simDevice) makeReading() core.Reading {
	elapsed := time.Since(s.start)
	return core.Reading{
		Received: time.Now(),
		Temp1:    math.Sin(elapsed.Seconds()/60)*100 + 100 + s.calibration.Temp1,
		Temp2:    math.Cos(elapsed.Seconds()/60)*100 + 100 + s.calibration.Temp2,
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

func (s *simDevice) GetCalibration() core.Calibration {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.calibration
}

func (s *simDevice) SetCalibration(calibration core.Calibration) {
	s.lock.Lock()
	s.calibration = calibration
	s.lock.Unlock()
	s.log.Info("calibration is now ", calibration)
}
