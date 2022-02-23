package client

import (
	"github.com/sculas/moonlight/network/decoder"
	"github.com/sculas/moonlight/util"
)

func (c *Client) StartReceiving() {
	for {
		frame, have := <-c.Receiver
		if !have || util.InvalidFrame(frame) {
			break
		}

		c.rb.Write(frame)
		packet, err := decoder.PacketDecoder(c.rb)
		if err != nil {
			c.log.Debugf("Error decoding packet: %s", err)
			c.Close()
			return
		}
		c.log.Debugf("got packet: %v", packet)

		// we're done
		c.ResetBuffers()
	}

	c.log.Debug("packet handler goroutine stopped")
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
	_ = c.c.Close(nil)
}
