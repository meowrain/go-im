package service

import (
	"testing"
)

func TestGetUser(t *testing.T) {
	user := UserService{}
	register, err := user.Register("17536720210", "test1234", "test@test.com", "1234567890", "1")
	if err != nil {
		t.Log(err)
	}
	t.Log(register)
}
