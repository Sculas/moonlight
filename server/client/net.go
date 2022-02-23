package client

import (
	"fmt"
	"github.com/sculas/moonlight/network/decoder"
	"github.com/sculas/moonlight/network/pk"
	"github.com/sculas/moonlight/util"
)

func (c *Client) WritePacket(packet pk.Packet) {
	if c.c == nil {
		return
	}

	if err := decoder.PacketEncoder(packet, c.c); err != nil {
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

		for i := 0; i < 2; i++ {
			fmt.Println("running packet stuff thingy")

			packet, _ /* safe to continue */, err := decoder.PacketDecoder(c.rb)
			if err != nil {
				c.Log.Debugf("Error decoding packet: %s", err)
				/*if !stc { // TODO better error handling
					c.Close()
					return
				}*/
			} else {
				if h, ok := GetHandler(packet.ID()); ok {
					if err = h.Handle(packet, c); err != nil {
						c.Log.Errorf("Error handling packet: %s", err)
						c.Close()
						return
					}
				} else {
					c.Log.Debugf("Packet %d has no handler", packet.ID())
				}
			}
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
