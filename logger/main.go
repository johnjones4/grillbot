package main

import (
	"context"
	"flag"
	"main/core"
	"main/device"
	"main/session"
	"main/ui"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	method := flag.String("method", "smoked", "Method of cooking")
	food := flag.String("food", "", "Food being prepared")
	changeThreshold := flag.Float64("change-threshold", 0.01, "Percent change threshold")
	timeThreshold := flag.Duration("time-threshold", time.Second*30, "Time threshold")
	simulated := flag.Bool("simulated", false, "use simulated data")
	resume := flag.String("resume", "", "Resume a previous cook")

	flag.Parse()

	ctx := context.Background()

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.Info("Initializing")

	var err error
	var sess core.Session
	if *resume == "" {
		md := core.Metadata{
			Food:   *food,
			Method: *method,
			Start:  time.Now(),
		}
		sess, err = session.New(log, md)
		if err != nil {
			log.Panic(err)
		}
	} else {
		sess, err = session.Open(log, *resume)
		if err != nil {
			log.Panic(err)
		}
	}

	log.Info("Session initialized")

	var dev core.Device
	if *simulated {
		dev = device.NewSim(log, sess)
	} else {
		dev, err = device.New(log, sess, *changeThreshold, *timeThreshold)
		if err != nil {
			log.Panic(err)
		}
	}
	log.Info("Device initialized")

	log.Info("Starting device")
	errChan := make(chan error)
	go dev.Start(ctx, errChan)

	uiinst, err := ui.New(log, sess, dev)
	if err != nil {
		log.Panic(err)
	}
	log.SetOutput(uiinst.LogView)
	sess.AddListener(uiinst.Listener())
	log.Info("UI starting")

	uiinst.Start(ctx, errChan)
}
