package client

import (
	"errors"
	"fmt"
	"github.com/sculas/moonlight/network/pk"
	"github.com/sculas/moonlight/server/client/state"
	"github.com/sirupsen/logrus"
)

type Handshake struct{}

func init() {
	RegisterHandler(pk.IDHandshake, &Handshake{})
}

func (h *Handshake) Handle(gp pk.Packet, c *Client) error {
	p := gp.(*pk.Handshake)
	c.Log.WithFields(logrus.Fields{
		"version":    p.ProtocolVersion,
		"address":    p.ServerAddress,
		"port":       p.ServerPort,
		"next_state": p.NextState,
	}).Debug("Received handshake!")
	s := state.From(p.NextState)
	if state.Invalid(s) {
		return errors.New(fmt.Sprintf("Invalid state: %d", p.NextState))
	}
	c.SetState(s)
	return nil
}
