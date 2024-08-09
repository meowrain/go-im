package meowlog

import "testing"

func TestLogger(t *testing.T) {
	logger := NewLogger("console", "debug", "logs")
	logger.Debug("test debug")
}
