package users_test

import "testing"

func TestLoginHandler(t *testing.T) {
	t.Run("Login with valid userId", func(t *testing.T) {
		// pass
	})
	t.Run("Login with missing userId", func(t *testing.T) {
		// pass
	})
	t.Run("Login without db connection", func(t *testing.T) {
		// pass
	})
	t.Run("Login with closed db connection", func(t *testing.T) {
		// pass
	})
}
