package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	CFG *Config
	LOG *log.Logger
	SRV *Server
	LOC *time.Location
)

func initInterrupt() {
	LOG.Println("-- start --")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(c chan os.Signal) {
		for range c {
			SRV.Shutdown()
			LOG.Println("-- stop --")
			os.Exit(137)
		}
	}(c)
}

func init() {
	CFG = ConfigInit()
	LOG = initLogging()
	initInterrupt()
	LOG.Printf("config path: %s", CFG.Path)
	var err error
	if LOC, err = time.LoadLocation(CFG.Service.Location); err != nil {
		eh(err)
	}
	SRV = new(Server)
}
