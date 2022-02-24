package client

import (
	"github.com/sculas/moonlight/network/pk"
	"github.com/sculas/moonlight/server/client/state"
)

type Handler interface {
	Handle(gp pk.Packet, c *Client) error
}

var handlers = map[state.ClientState]map[uint32]Handler{
	state.Handshaking: {},
	state.Play:        {},
	state.Status:      {},
	state.Login:       {},
}

func RegisterHandler(state state.ClientState, id uint32, h Handler) {
	handlers[state][id] = h
}

func GetHandler(state state.ClientState, id uint32) (Handler, bool) {
	h, ok := handlers[state][id]
	return h, ok
}

func (c *Client) handlePacket(state state.ClientState, packet pk.Packet) bool {
	if h, ok := GetHandler(state, packet.ID()); ok {
		if err := h.Handle(packet, c); err != nil {
			c.Log.Errorf("error handling packet: %s", err)
			c.Close()
			return false
		}
	} else {
		c.Log.Warnf("packet %T has no handler", packet)
	}
	return true
}
