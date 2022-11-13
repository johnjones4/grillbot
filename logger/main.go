package main

import (
	"main/device"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	dev, err := device.New(log)
	if err != nil {
		log.Panic(err)
	}

	dev.AddListener(func(d *device.Device, m device.Message) {
		log.Info(m)
	})

	for true {
	}
}
