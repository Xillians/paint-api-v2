package paints_test

import (
	"paint-api/internal/handlers/paints"
	"testing"
)

func TestValidateColorCode(t *testing.T) {
	tests := []struct {
		colorCode string
		expected  bool
	}{
		{"#FFFFFF", true},
		{"#000000", true},
		{"#FF5733", true},
		{"#123456", true},
		{"#ABCDEF", true},
		{"#GHIJKL", false},
		{"123456", false},
		{"#12345G", false},
	}

	for _, test := range tests {
		result := paints.ValidateColorCode(test.colorCode)
		if result != test.expected {
			t.Errorf("Expected %v for color code %s, got %v", test.expected, test.colorCode, result)
		}
	}
}
