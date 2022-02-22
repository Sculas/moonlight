package server

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/config"
	"github.com/sculas/moonlight/global"
	"github.com/sculas/moonlight/server/client"
	"github.com/sirupsen/logrus"
)

// TODO:
//  Need to decide on whether we should do this:
//    it would also be good if every connection has a buffer of their own, so we don't have to get one from the pool for every packet

type Server struct {
	*gnet.BuiltinEventEngine

	log *logrus.Entry
}

func New() *Server {
	return &Server{
		log: global.ServerLogger,
	}
}

func (s *Server) OnBoot(gnet.Engine) gnet.Action {
	s.log.Infof("Listening on port %d.", config.Config.Server.Port)
	return gnet.None
}

// OnShutdown fires when the server is being shut down, it is called right after
// all event-loops and connections are closed.
func (s *Server) OnShutdown(gnet.Engine) {
	s.log.Info("Shutting down.")
}

// OnOpen fires when a new connection has been opened.
// The parameter out is the return value which is going to be sent back to the peer.
func (s *Server) OnOpen(conn gnet.Conn) (out []byte, action gnet.Action) {
	s.log.Debugf("new connection from %s", conn.RemoteAddr())

	c := client.NewClient(conn)
	conn.SetContext(c)

	// start receiver
	go c.StartReceiving()

	return
}

// OnClose fires when a connection has been closed.
// The parameter err is the last known connection error.
func (s *Server) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	s.log.Debugf("connection from %s closed: %s", conn.RemoteAddr(), err)

	c := conn.Context().(*client.Client)
	c.Cleanup()

	return
}

// OnTraffic fires when a local socket receives data from the peer.
func (s *Server) OnTraffic(conn gnet.Conn) gnet.Action {
	buf, _ := conn.Next(-1)
	s.log.Debugf("traffic from %s: %d bytes", conn.RemoteAddr(), len(buf))

	c := conn.Context().(*client.Client)
	c.Receiver <- buf[:] // pass copy of buffer to receiver

	return gnet.None
}
