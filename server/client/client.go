package client

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/global"
	"github.com/sculas/moonlight/network/serde"
	"github.com/sculas/moonlight/server/client/state"
	"github.com/sirupsen/logrus"
)

type Client struct {
	// connection
	c gnet.Conn

	// buffers
	rb *serde.ByteBuf

	// frame receiver
	Receiver chan []byte

	// logger
	log *logrus.Entry

	// client state
	state state.ClientState
}

func NewClient(c gnet.Conn) *Client {
	return &Client{
		c: c,

		rb: serde.NewByteBuf(),

		Receiver: make(chan []byte),

		log: global.ClientLogger.WithField("addr", c.RemoteAddr().String()), // TODO
	}
}
