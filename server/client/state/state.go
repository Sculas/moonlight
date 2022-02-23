package state

type ClientState int

const (
	invalid ClientState = iota - 2

	Handshaking
	Play
	Status
	Login
)

func Invalid(state ClientState) bool {
	return state == invalid
}

func From(state int) ClientState {
	switch state {
	case -1:
		return Handshaking
	case 0:
		return Play
	case 1:
		return Status
	case 2:
		return Login
	default:
		return invalid
	}
}
