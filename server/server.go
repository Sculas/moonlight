package server

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/config"
	"github.com/sirupsen/logrus"
)

type server struct {
	*gnet.BuiltinEventEngine

	log *logrus.Entry
}

func New(logger *logrus.Entry) *server {
	return &server{
		log: logger,
	}
}

func (s *server) OnBoot(gnet.Engine) gnet.Action {
	s.log.Infof("Listening on port %d.", config.Config.Server.Port)
	return gnet.None
}

func (s *server) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}
