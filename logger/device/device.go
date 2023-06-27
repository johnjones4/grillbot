package device

import (
	"context"
	"strings"
	"time"

	"main/core"

	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type device struct {
	baseDevice
	port   *serial.Port
	buffer []byte
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
		buffer: make([]byte, 0),
	}

	return d, nil
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
	d.baseDevice.start(ctx, outChan, d.port.Close, d.nextMessage)
}
