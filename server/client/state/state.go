package state

type ClientState int

const (
	Handshaking ClientState = iota - 1
	Play
	Status
	Login
)
