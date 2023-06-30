package main

import (
	"context"
	"flag"
	"main/core"
	"main/device"
	"main/server"
	"main/session"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	method := flag.String("method", "smoked", "Method of cooking")
	food := flag.String("food", "", "Food being prepared")
	changeThreshold := flag.Float64("change-threshold", 0.02, "Percent change threshold")
	timeThreshold := flag.Duration("time-threshold", time.Second*30, "Time threshold")
	simulated := flag.Bool("simulated", false, "use simulated data")
	file := flag.String("file", "", "Resume a previous cook")
	serial := flag.String("serial", "", "Serial device to use")
	debug := flag.Bool("debug", false, "Run with debug logging")
	host := flag.String("host", ":8080", "Hostname to listen on")

	flag.Parse()

	ctx := context.Background()

	log := logrus.New()
	if *debug {
		log.SetLevel(logrus.DebugLevel)
	}
	log.Info("Initializing")

	var err error
	var sess core.Session
	if *file == "" {
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
		sess, err = session.Open(log, *file)
		if err != nil {
			log.Panic(err)
		}
	}

	log.Info("Session initialized")

	var dev core.Device
	if *simulated {
		dev = device.NewSim(log, sess)
	} else {
		dev, err = device.New(log, sess, *changeThreshold, *timeThreshold, *serial)
		if err != nil {
			log.Panic(err)
		}
	}
	log.Info("Device initialized")

	log.Info("Starting device")
	go dev.Start(ctx)

	api := server.New(sess, dev, log)
	err = http.ListenAndServe(*host, api.Mux)
	if err != nil {
		log.Panic(err)
	}
}
