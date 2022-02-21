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

// OnShutdown fires when the server is being shut down, it is called right after
// all event-loops and connections are closed.
func (s *server) OnShutdown(gnet.Engine) {
	s.log.Info("Shutting down.")
}

// OnOpen fires when a new connection has been opened.
// The parameter out is the return value which is going to be sent back to the peer.
func (s *server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	s.log.Debugf("new connection from %s", c.RemoteAddr())
	return
}

// OnClose fires when a connection has been closed.
// The parameter err is the last known connection error.
func (s *server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	s.log.Debugf("connection from %s closed: %s", c.RemoteAddr(), err)
	return
}

// OnTraffic fires when a local socket receives data from the peer.
func (s *server) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	s.log.Debugf("recv %d bytes\n", len(buf))
	c.Write(buf)
	return gnet.None
}
