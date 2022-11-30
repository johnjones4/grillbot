package device

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
	"time"

	"main/core"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/darwin"
	"github.com/sirupsen/logrus"
)

var tempCharacteristicUUID = ble.MustParse("a3612fbb-7c00-4ab2-b925-425c4ef2a002")
var calibCharacteristicUUID = ble.MustParse("09222388-fd96-4194-822b-fa052786c130")

type device struct {
	client              ble.Client
	tempCharacteristic  *ble.Characteristic
	calibCharacteristic *ble.Characteristic
	log                 *logrus.Logger
	lastValue           []byte
	lastMessage         core.Reading
	lastUpdate          time.Time
	session             core.Session
	deltaThreshold      float64
	timeThreshold       time.Duration
	currentCaliration   core.Calibration
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

	tempCharacteristicI := profile.Find(ble.NewCharacteristic(tempCharacteristicUUID))
	if tempCharacteristicI == nil {
		return nil, errors.New("temp characteristic not found")
	}
	tempCharacteristic := tempCharacteristicI.(*ble.Characteristic)
	log.Debug("Got temp characteristic")

	calibCharacteristicI := profile.Find(ble.NewCharacteristic(tempCharacteristicUUID))
	if calibCharacteristicI == nil {
		return nil, errors.New("temp characteristic not found")
	}
	calibCharacteristic := calibCharacteristicI.(*ble.Characteristic)
	log.Debug("Got temp characteristic")

	d := &device{
		client:              client,
		tempCharacteristic:  tempCharacteristic,
		calibCharacteristic: calibCharacteristic,
		log:                 log,
		lastValue:           nil,
		session:             sess,
		deltaThreshold:      deltaThreshold,
		timeThreshold:       timeThreshold,
	}

	val, err := client.ReadCharacteristic(tempCharacteristic)
	if err != nil {
		return nil, err
	}
	d.lastValue = val
	log.Debug("Initial temp value: ", d.lastValue)

	m, err := parseMessage(val)
	if err != nil {
		return nil, err
	}
	d.lastMessage = m
	log.Debug("Initial message: ", d.lastMessage)

	calibVal, err := client.ReadCharacteristic(calibCharacteristic)
	if err != nil {
		return nil, err
	}
	log.Debug("Initial calbiration value: ", d.lastValue)
	calibration, err := parseCalibration(calibVal)
	if err != nil {
		return nil, err
	}
	d.currentCaliration = calibration
	log.Debug("Initial calibration: ", calibration)

	if isValidReading(m) {
		d.session.NewReading(m)
	}

	return d, nil
}

func filter(a ble.Advertisement) bool {
	return a.LocalName() == "GrillBot"
}

func isValidReading(r core.Reading) bool {
	return r.Temp1 != 0 && r.Temp2 != 0 && !r.Received.IsZero()
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

func parseCalibration(b []byte) (core.Calibration, error) {
	if len(b) == 0 {
		return core.Calibration{}, nil
	}
	return core.Calibration{
		Temp1: math.Float64frombits(binary.LittleEndian.Uint64(b[:8])),
		Temp2: math.Float64frombits(binary.LittleEndian.Uint64(b[8:])),
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
			val, err := d.client.ReadCharacteristic(d.tempCharacteristic)
			if err != nil {
				outChan <- err
				return
			}
			if bytes.Equal(val, d.lastValue) {
				continue
			}
			d.log.Debug("New data available: ", hex.EncodeToString(val))
			m, err := parseMessage(val)
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

func (d *device) GetCalibration() core.Calibration {
	return d.currentCaliration
}

func (d *device) SetCalibration(c core.Calibration) {
	d.log.Info("Sending new calibratin data: ", c)
	buf1 := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf1, math.Float64bits(c.Temp1))
	buf2 := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf2, math.Float64bits(c.Temp2))
	buf3 := append(buf1, buf2...)
	d.log.Info("Encoded calibration data as: ", hex.EncodeToString(buf3))
	d.calibCharacteristic.SetValue(buf3)
}
