package client

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/sculas/moonlight/network/pipeline"
)

type Client struct {
	c        *gnet.Conn
	pipeline *pipeline.ChannelPipeline
}

func NewClient(c *gnet.Conn) *Client {
	client := &Client{c: c}
	client.pipeline = pipeline.New(client)
	return client
}

func (c *Client) Receive(conn gnet.Conn) {
	c.pipeline.Fire(conn)
}
