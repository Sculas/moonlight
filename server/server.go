package server

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/config"
	"github.com/sculas/moonlight/server/client"
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

	// FIXME:
	//  let's implement something like Netty's MessageToByte and ByteToMessage (en/de)coders.
	//  it would also be good if every connection has a buffer of their own, so we don't have to get one from the pool for every packet

	cc := client.NewClient(&c) // FIXME: are we allowed to hold a reference of the conn?
	c.SetContext(cc)

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
	//c.Context().(*pipeline.ChannelPipeline).Fire(c)
	return gnet.None
}
