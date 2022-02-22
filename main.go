package main

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/config"
	"github.com/sculas/moonlight/global"
	"github.com/sculas/moonlight/server"
	"github.com/sculas/moonlight/util"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	config.Initialize()

	global.Logger = logrus.New()
	global.Logger.SetFormatter(&nested.Formatter{
		HideKeys:      true,
		ShowFullLevel: true,
	})
	global.Logger.Level = logrus.DebugLevel

	global.ServerLogger = global.Logger.WithField(util.Component("server"))
	global.ClientLogger = global.Logger.WithField(util.Component("client"))

	global.Server = server.New()
	log.Fatal(gnet.Run(
		global.Server,
		fmt.Sprintf("tcp://:%d", config.Config.Server.Port),
		gnet.WithOptions(gnet.Options{
			Logger: global.Logger.WithField(util.Component("gnet")),

			Multicore: config.Config.Server.Multicore,
			ReuseAddr: config.Config.Server.ReuseAddr,
			ReusePort: config.Config.Server.ReusePort,
		}),
	))
}
