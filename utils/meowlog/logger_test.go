package meowlog

import "testing"

func TestLogger(t *testing.T) {
	logger := NewLogger("file", "debug", "logs")
	logger.Debug("test debug")
}
