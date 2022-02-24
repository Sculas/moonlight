package handler

import (
	"github.com/sculas/moonlight/network/pk"
	"github.com/sculas/moonlight/server/client"
	"github.com/sculas/moonlight/server/client/state"
)

type StatusInStart struct{}

func init() {
	client.RegisterHandler(state.Status, pk.IDStatusInStart, &StatusInStart{})
}

func (StatusInStart) Handle(_ pk.Packet, c *client.Client) error {
	c.Log.Debug("Received StatusInStart!")
	return nil
}
