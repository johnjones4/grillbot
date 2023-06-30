package device

import (
	"context"
	_ "embed"
	"main/core"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

//go:embed sampledata.txt
var sampledata string

type simDevice struct {
	baseDevice
	lock  sync.RWMutex
	start time.Time
	data  [][]byte
}

func NewSim(log *logrus.Logger, sess core.Session) core.Device {
	lines := strings.Split(sampledata, "\n")
	buffer := make([][]byte, len(lines))
	for i := 0; i < len(lines); i++ {
		buffer[i] = []byte(lines[i])
	}
	return &simDevice{
		baseDevice: baseDevice{
			log:     log,
			session: sess,
		},
		lock:  sync.RWMutex{},
		start: time.Now(),
		data:  buffer,
	}
}

func (s *simDevice) Start(ctx context.Context) {
	s.baseDevice.start(ctx, func() error { return nil }, s.nextMessage)
}

func (d *simDevice) nextMessage() ([]byte, error) {
	time.Sleep(time.Second)
	elapsed := time.Since(d.start)
	return d.data[int(elapsed)%len(d.data)], nil
}
