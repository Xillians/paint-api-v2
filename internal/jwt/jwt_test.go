package jwt_test

import (
	jwtService "paint-api/internal/jwt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestJwtClient(t *testing.T) {
	client := jwtService.NewJWTService("my_secret")

	t.Run("Generate a token", func(t *testing.T) {
		_, err := client.GenerateToken("1234567890", "administrator")
		if err != nil {
			t.Errorf("Error generating token: %v", err)
		}
	})
	t.Run("Verify a token", func(t *testing.T) {
		token, err := client.GenerateToken("1234567890", "administrator")
		if err != nil {
			t.Errorf("Error generating token: %v", err)
		}

		output, err := client.VerifyToken(token)
		if err != nil {
			t.Errorf("Error verifying token: %v", err)
		}
		claims := output.Claims.(jwt.MapClaims)
		if claims["user_id"] != "1234567890" {
			t.Errorf("Expected user_id to be 1234567890, got %v", claims["user_id"])
		}
		if claims["role"] != "administrator" {
			t.Errorf("Expected role to be administrator, got %v", claims["role"])
		}
	})
	t.Run("Send invalid token", func(t *testing.T) {
		_, err := client.GenerateToken("1234567890", "administrator")
		if err != nil {
			t.Errorf("Error generating token: %v", err)
		}

		// send in a token with a different signature
		fakeClient := jwtService.NewJWTService("fake_secret")
		fakeToken, err := fakeClient.GenerateToken("1234567890", "administrator")
		if err != nil {
			t.Errorf("Error generating token: %v", err)
		}

		_, err = client.VerifyToken(fakeToken)
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}
