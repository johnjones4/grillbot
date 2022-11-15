package main

import (
	"context"
	"main/core"
	"main/device"
	"main/outputs"
	"main/session"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	md := core.Metadata{
		Food:   "brisket",
		Method: "smoked",
		Start:  time.Now(),
	}
	sess, err := session.New(log, md)
	if err != nil {
		log.Panic(err)
	}

	dev, err := device.New(log, sess)
	if err != nil {
		log.Panic(err)
	}

	termOut, err := outputs.NewTable()
	if err != nil {
		log.Panic(err)
	}
	sess.AddListener(termOut.Listener())
	go termOut.Start(ctx)

	serverOut := outputs.NewServer(sess, ":8080")
	sess.AddListener(serverOut.Listener())
	go serverOut.Start(ctx)

	dev.Start(ctx)
}
