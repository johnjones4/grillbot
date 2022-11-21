package main

import (
	"context"
	"flag"
	"main/core"
	"main/device"
	"main/outputs"
	"main/session"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	method := flag.String("method", "smoked", "Method of cooking")
	food := flag.String("food", "", "Food being prepared")
	httpHost := flag.String("host", ":8080", "HTTP host string")
	changeThreshold := flag.Float64("change-threshold", 0.01, "Percent change threshold")
	timeThreshold := flag.Duration("time-threshold", time.Second*30, "Time threshold")

	flag.Parse()

	ctx := context.Background()

	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	md := core.Metadata{
		Food:   *food,
		Method: *method,
		Start:  time.Now(),
	}
	sess, err := session.New(log, md)
	if err != nil {
		log.Panic(err)
	}

	dev, err := device.New(log, sess, *changeThreshold, *timeThreshold)
	if err != nil {
		log.Panic(err)
	}

	errChan := make(chan error)

	termOut, err := outputs.NewTable()
	if err != nil {
		log.Panic(err)
	}
	sess.AddListener(termOut.Listener())
	go termOut.Start(ctx, errChan)

	serverOut := outputs.NewServer(sess, *httpHost)
	sess.AddListener(serverOut.Listener())
	go serverOut.Start(ctx, errChan)

	go dev.Start(ctx, errChan)

	for err := range errChan {
		if err != nil {
			log.Panic(err)
		}
		os.Exit(0)
	}
}
