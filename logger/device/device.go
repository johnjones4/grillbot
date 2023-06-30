package device

import (
	"bufio"
	"context"
	"time"

	"main/core"

	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type device struct {
	baseDevice
	port   *serial.Port
	buffer *bufio.Reader
}

func New(log *logrus.Logger, sess core.Session, deltaThreshold float64, timeThreshold time.Duration, handle string) (core.Device, error) {
	c := &serial.Config{Name: handle, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	d := &device{
		baseDevice: baseDevice{
			log:            log,
			lastValue:      nil,
			session:        sess,
			deltaThreshold: deltaThreshold,
			timeThreshold:  timeThreshold,
		},
		port:   s,
		buffer: bufio.NewReader(s),
	}

	return d, nil
}

func (d *device) nextMessage() ([]byte, error) {
	return d.buffer.ReadBytes('\n')
}

func (d *device) Start(ctx context.Context) {
	d.baseDevice.start(ctx, d.port.Close, d.nextMessage)
}
