package unit

import (
	"testing"
)

// Example unit test - replace with your actual tests
func TestExample(t *testing.T) {
	// Arrange
	expected := 4

	// Act
	result := 2 + 2

	// Assert
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

// Example table-driven test
func TestAddition(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed numbers", -2, 3, 1},
		{"zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a + tt.b
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
