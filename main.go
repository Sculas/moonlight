package main

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/config"
	"github.com/sculas/moonlight/server"
	"github.com/sculas/moonlight/util"
	"github.com/sirupsen/logrus"
)

func main() {
	config.Initialize()

	log := logrus.New()
	log.SetFormatter(&nested.Formatter{
		HideKeys: true,
	})
	log.Level = logrus.DebugLevel

	srv := server.New(log.WithField(util.Component("server")))
	log.Fatal(gnet.Run(
		srv,
		fmt.Sprintf("tcp://:%d", config.Config.Server.Port),
		gnet.WithOptions(gnet.Options{
			Logger: log.WithField(util.Component("gnet")),

			Multicore: config.Config.Server.Multicore,
			ReuseAddr: config.Config.Server.ReuseAddr,
			ReusePort: config.Config.Server.ReusePort,
		}),
	))
}
