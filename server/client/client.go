package client

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/global"
	"github.com/sculas/moonlight/network/serde"
	"github.com/sculas/moonlight/util"
	"github.com/sirupsen/logrus"
)

type Client struct {
	// connection
	c gnet.Conn

	// buffers
	rb, wb *serde.ByteBuf

	// frame receiver
	Receiver chan []byte

	// logger
	log *logrus.Entry
}

func NewClient(c gnet.Conn) *Client {
	return &Client{
		c: c,

		rb: serde.NewByteBuf(),
		wb: serde.NewByteBuf(),

		Receiver: make(chan []byte),

		log: global.ClientLogger.WithField("addr", c.RemoteAddr().String()), // TODO
	}
}

func (c *Client) StartReceiving() {
	for {
		frame, have := <-c.Receiver
		if !have || util.InvalidFrame(frame) {
			break
		}

		c.log.Debugf("got traffic in our goroutine: %s", string(frame))
	}

	c.log.Debug("packet handler goroutine stopped")
	c.Close()
}

func (c *Client) StopReceiving() {
	close(c.Receiver)
}

func (c *Client) ResetBuffers() {
	c.rb.Reset()
	c.wb.Reset()
}

func (c *Client) Cleanup() {
	c.StopReceiving()
	c.ResetBuffers()
	c.Close()
}

func (c *Client) Close() {
	_ = c.c.Close(nil)
}
