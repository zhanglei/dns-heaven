package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/greenboxal/dns-heaven"
	"github.com/greenboxal/dns-heaven/osx"
	"github.com/sirupsen/logrus"
)

var config = &dnsheaven.Config{}

func init() {
	flag.StringVar(&config.Address, "address", "127.0.0.1:53", "address to listen")
	flag.IntVar(&config.Timeout, "timeout", 2000, "request timeout")
	flag.IntVar(&config.Interval, "interval", 1000, "interval between requests")
}

func main() {
	flag.Parse()

	resolver, err := osx.New(config)

	if err != nil {
		logrus.WithError(err).Error("error starting server")
		os.Exit(1)
	}

	server := dnsheaven.NewServer(config, resolver)

	stopping := false
	go func() {
		err := server.Start()

		if !stopping && err != nil {
			logrus.WithError(err).Error("error starting server")
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	_ = <-sig

	server.Shutdown()
}
