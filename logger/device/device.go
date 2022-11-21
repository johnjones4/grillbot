package device

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"math"
	"time"

	"main/core"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/darwin"
	"github.com/sirupsen/logrus"
)

var characteristicUUID = ble.MustParse("a3612fbb-7c00-4ab2-b925-425c4ef2a002")

type device struct {
	client         ble.Client
	characteristic *ble.Characteristic
	log            *logrus.Logger
	lastValue      []byte
	lastMessage    core.Reading
	lastUpdate     time.Time
	session        core.Session
	deltaThreshold float64
	timeThreshold  time.Duration
}

func New(log *logrus.Logger, sess core.Session, deltaThreshold float64, timeThreshold time.Duration) (core.Device, error) {
	dev, err := darwin.NewDevice()
	if err != nil {
		return nil, err
	}
	ble.SetDefaultDevice(dev)
	log.Debug("Loaded device")

	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 60*time.Second))
	client, err := ble.Connect(ctx, filter)
	if err != nil {
		return nil, err
	}
	log.Debug("Got client")

	profile, err := client.DiscoverProfile(true)
	if err != nil {
		return nil, err
	}
	log.Debug("Got profile")

	characteristicI := profile.Find(ble.NewCharacteristic(characteristicUUID))
	if characteristicI == nil {
		return nil, errors.New("characteristic not found")
	}
	characteristic := characteristicI.(*ble.Characteristic)
	log.Debug("Got characteristic")

	d := &device{
		client:         client,
		characteristic: characteristic,
		log:            log,
		lastValue:      nil,
		session:        sess,
		deltaThreshold: deltaThreshold,
		timeThreshold:  timeThreshold,
	}

	val, err := client.ReadCharacteristic(characteristic)
	if err != nil {
		return nil, err
	}
	d.lastValue = val
	log.Debug("Initial value: ", d.lastValue)

	m, err := parseMessage(val)
	if err != nil {
		return nil, err
	}
	d.lastMessage = m
	d.session.NewReading(m)
	log.Debug("Initial message: ", d.lastMessage)

	return d, nil
}

func filter(a ble.Advertisement) bool {
	return a.LocalName() == "GrillBot"
}

func parseMessage(b []byte) (core.Reading, error) {
	if len(b) == 0 {
		return core.Reading{}, nil
	}
	return core.Reading{
		Received: time.Now(),
		Temp1:    math.Float64frombits(binary.LittleEndian.Uint64(b[:8])),
		Temp2:    math.Float64frombits(binary.LittleEndian.Uint64(b[8:])),
	}, nil
}

func (d *device) Start(ctx context.Context, outChan chan error) {
	stop := ctx.Done()
	timer := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-stop:
			return
		case <-timer.C:
			val, err := d.client.ReadCharacteristic(d.characteristic)
			if err != nil {
				outChan <- err
				return
			}
			if bytes.Equal(val, d.lastValue) {
				continue
			}
			d.log.Debug("New data available")
			m, err := parseMessage(val)
			if err != nil {
				outChan <- err
				return
			}
			d.log.Debug("New message ", m)
			diff := d.lastMessage.MaxPcntDifference(m)
			now := time.Now()
			if (diff <= d.deltaThreshold && now.Add(d.timeThreshold*-1).Before(d.lastUpdate)) || d.lastMessage.Received == m.Received {
				continue
			}
			d.lastUpdate = now
			d.lastValue = val
			d.lastMessage = m
			d.session.NewReading(m)
		}
	}
}
