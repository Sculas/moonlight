package util

func InvalidFrame(frame []byte) bool {
	return frame == nil || len(frame) == 0
}
