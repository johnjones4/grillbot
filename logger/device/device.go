package device

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"math"
	"strings"
	"sync"
	"time"

	"main/core"

	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type device struct {
	log            *logrus.Logger
	lastValue      []byte
	lastMessage    core.Reading
	lastUpdate     time.Time
	session        core.Session
	deltaThreshold float64
	timeThreshold  time.Duration
	calibration    [2]float64
	port           *serial.Port
	buffer         []byte
	lock           sync.RWMutex
}

func New(log *logrus.Logger, sess core.Session, deltaThreshold float64, timeThreshold time.Duration, handle string) (core.Device, error) {
	c := &serial.Config{Name: handle, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	d := &device{
		log:            log,
		lastValue:      nil,
		session:        sess,
		deltaThreshold: deltaThreshold,
		timeThreshold:  timeThreshold,
		port:           s,
		buffer:         make([]byte, 0),
	}

	return d, nil
}

func isValidReading(r core.Reading) bool {
	return r.Temperatures[0] != 0 && r.Temperatures[1] != 0 && !r.Received.IsZero()
}

func (d *device) parseMessage(b []byte) (core.Reading, error) {
	if len(b) == 0 || len(b) < 16 {
		return core.Reading{}, nil
	}
	return core.Reading{
		Received: time.Now(),
		Temperatures: [2]float64{
			math.Float64frombits(binary.LittleEndian.Uint64(b[:8])) + d.calibration[0],
			math.Float64frombits(binary.LittleEndian.Uint64(b[8:16])) + d.calibration[1],
		},
	}, nil
}

func (d *device) nextMessage() ([]byte, error) {
	buffer := make([]byte, 17)
	read, err := d.port.Read(buffer)
	if err != nil {
		return nil, err
	}

	buffer = append(d.buffer, buffer[:read]...)

	newline := strings.IndexByte(string(buffer), '\n')

	if newline < 0 {
		d.buffer = buffer
		return nil, nil
	}

	d.buffer = buffer[newline:]
	return buffer[:newline], nil
}

func (d *device) Start(ctx context.Context, outChan chan error) {
	stop := ctx.Done()
	timer := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-stop:
			err := d.port.Close()
			outChan <- err
			return
		case <-timer.C:
			val, err := d.nextMessage()
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

func (d *device) GetCalibration() [2]float64 {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.calibration
}

func (d *device) SetCalibration(calibration [2]float64) {
	d.lock.Lock()
	d.calibration = calibration
	d.lock.Unlock()
	d.log.Info("calibration is now ", calibration)
}
