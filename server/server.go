package server

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/config"
	"github.com/sculas/moonlight/global"
	"github.com/sculas/moonlight/server/client"
	"github.com/sirupsen/logrus"
)

var Server *server

type server struct {
	*gnet.BuiltinEventEngine

	log *logrus.Entry
}

func New() *server {
	return &server{
		log: global.ServerLogger,
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
func (s *server) OnOpen(conn gnet.Conn) (out []byte, action gnet.Action) {
	s.log.Debugf("new connection from %s", conn.RemoteAddr())

	c := client.NewClient(conn)
	conn.SetContext(c)

	// start receiver
	go c.StartReceiving()

	return
}

// OnClose fires when a connection has been closed.
// The parameter err is the last known connection error.
func (s *server) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	s.log.Debugf("connection from %s closed: %s", conn.RemoteAddr(), err)

	c := conn.Context().(*client.Client)
	c.Cleanup()

	return
}

// OnTraffic fires when a local socket receives data from the peer.
func (s *server) OnTraffic(conn gnet.Conn) gnet.Action {
	buf, _ := conn.Next(-1)
	s.log.Debugf("traffic from %s: %d bytes", conn.RemoteAddr(), len(buf))

	c := conn.Context().(*client.Client)
	c.Receiver <- buf[:] // pass copy of buffer to receiver

	return gnet.None
}
