package device

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"main/core"
	"math"
	"time"

	"github.com/sirupsen/logrus"
)

type baseDevice struct {
	log            *logrus.Logger
	lastValue      []byte
	lastMessage    core.Reading
	lastUpdate     time.Time
	deltaThreshold float64
	timeThreshold  time.Duration
	session        core.Session
}

func (d *baseDevice) parseMessage(b []byte) (core.Reading, error) {
	if len(b) == 0 || len(b) < 16 {
		return core.Reading{}, nil
	}
	return core.Reading{
		Received: time.Now(),
		Temperatures: [2]float64{
			math.Float64frombits(binary.LittleEndian.Uint64(b[:8])),
			math.Float64frombits(binary.LittleEndian.Uint64(b[8:16])),
		},
	}, nil
}

func (d *baseDevice) start(ctx context.Context, outChan chan error, closer func() error, nextMessage func() ([]byte, error)) {
	stop := ctx.Done()
	timer := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-stop:
			outChan <- closer()
			return
		case <-timer.C:
			val, err := nextMessage()
			if err != nil {
				outChan <- err
				return
			}
			if bytes.Equal(val, d.lastValue) {
				continue
			}
			d.log.Debug("New data available: ", hex.EncodeToString(val))
			m, err := d.parseMessage(val)
			if err != nil {
				outChan <- err
				return
			}
			d.log.Debug("New message: ", m)
			diff := d.lastMessage.MaxPcntDifference(m)
			now := time.Now()
			if (diff <= d.deltaThreshold && now.Add(d.timeThreshold*-1).Before(d.lastUpdate)) || d.lastMessage.Received == m.Received {
				continue
			}
			d.lastUpdate = now
			d.lastValue = val
			d.lastMessage = m
			if isValidReading(m) {
				d.session.NewReading(m)
			}
		}
	}
}

func isValidReading(r core.Reading) bool {
	return r.Temperatures[0] != 0 && r.Temperatures[1] != 0 && !r.Received.IsZero()
}
