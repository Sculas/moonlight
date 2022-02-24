package handler

import (
	"errors"
	"fmt"
	"github.com/sculas/moonlight/network/pk"
	"github.com/sculas/moonlight/server/client"
	"github.com/sculas/moonlight/server/client/state"
	"github.com/sirupsen/logrus"
)

type HandshakingInSetProtocol struct{}

func init() {
	client.RegisterHandler(state.Handshaking, pk.IDHandshakingInSetProtocol, &HandshakingInSetProtocol{})
}

func (HandshakingInSetProtocol) Handle(gp pk.Packet, c *client.Client) error {
	p := gp.(*pk.HandshakingInSetProtocol)
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
