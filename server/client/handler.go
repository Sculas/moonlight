package client

import (
	"github.com/sculas/moonlight/network/pk"
)

type Handler interface {
	Handle(gp pk.Packet, c *Client) error
}

var handlers = map[uint32]Handler{}

func RegisterHandler(id uint32, h Handler) {
	handlers[id] = h
}

func GetHandler(id uint32) (Handler, bool) {
	h, ok := handlers[id]
	return h, ok
}
