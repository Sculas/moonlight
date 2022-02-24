package client

import (
	"github.com/sculas/moonlight/network/codec"
	"github.com/sculas/moonlight/network/pk"
	"github.com/sculas/moonlight/util"
)

func (c *Client) WritePacket(packet pk.Packet) {
	if c.c == nil {
		return
	}

	if err := codec.PacketEncoder(packet, c.c); err != nil {
		c.Log.Errorf("Error encoding packet: %v", err)
		c.Close()
	}
}

func (c *Client) StartReceiving() {
	for {
		frame, have := <-c.Receiver
		if !have || util.InvalidFrame(frame) {
			break
		}

		c.rb.Write(frame)

		// TODO better error handling
		err := codec.Decode(c.rb, &c.state, c.handlePacket)
		if err != nil {
			c.Log.Debugf("error decoding packet: %s", err)
		}

		// we're done
		c.ResetBuffers()
	}

	c.Log.Debug("packet handler goroutine stopped")
	c.Close()
}

func (c *Client) StopReceiving() {
	close(c.Receiver)
}

func (c *Client) ResetBuffers() {
	c.rb.Reset()
}

func (c *Client) Cleanup() {
	c.StopReceiving()
	c.ResetBuffers()
	c.Close()
}

func (c *Client) Close() {
	if c.c == nil {
		return
	}
	_ = c.c.Close(nil)
}
