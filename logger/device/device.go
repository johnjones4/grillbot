package device

import (
	"context"
	"errors"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/darwin"
	"github.com/sirupsen/logrus"
)

var characteristicUUID = ble.MustParse("a3612fbb-7c00-4ab2-b925-425c4ef2a002")

type Listener func(*Device, Message)

type Message struct {
	Temp1 float64
	Temp2 float64
}

type Device struct {
	characteristic *ble.Characteristic
	log            *logrus.Logger
	lastValue      string
	listeners      []Listener
}

func New(log *logrus.Logger) (*Device, error) {
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

	d := &Device{
		characteristic: characteristic,
		log:            log,
		lastValue:      "",
		listeners:      []Listener{},
	}

	val, err := client.ReadCharacteristic(characteristic)
	if err != nil {
		return nil, err
	}
	d.lastValue = string(val)
	log.Debug("Initial value: ", d.lastValue)

	client.Subscribe(characteristic, false, d.valueUpdated)

	return d, nil
}

func (d *Device) valueUpdated(req []byte) {
	d.lastValue = string(req)
	d.log.Info("Latest reading: %s", d.lastValue)
	for _, l := range d.listeners {
		l(d, Message{}) //TODO
	}
}

func (d *Device) AddListener(l Listener) {
	d.listeners = append(d.listeners, l)
}

func (d *Device) LastValue() string {
	return d.lastValue
}

func filter(a ble.Advertisement) bool {
	return a.LocalName() == "GrillBot"
}
