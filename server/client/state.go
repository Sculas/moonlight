package client

import "github.com/sculas/moonlight/server/client/state"

func (c *Client) State() state.ClientState {
	return c.state
}

func (c *Client) SetState(state state.ClientState) {
	c.state = state
}
