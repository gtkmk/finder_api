package tests

import (
	"testing"

	phone2 "github.com/gtkmk/finder_api/core/domain/phone"
)

func TestPhoneValidation(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{"123456789", true},        // Too short phone number
		{"123456789012", true},     // Too long phone number
		{"123abc4567", true},       // Non-numeric characters
		{"(31) 40028922", false},   // Valid masked phone number
		{"(31) 98598-8922", false}, // Valid masked cellphone number
		{"(31) 985988922", false},  // Valid masked cellphone number
		{"3140028922", false},      // Valid unmasked phone number
	}

	for _, test := range tests {
		phone := phone2.NewPhone(test.input)
		err := phone.Validate()

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input %s, but got nil", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
			}
		}
	}
}

func TestUnmaskedString(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"(31) 4002-8922", "3140028922", false},
		{"3140028922", "3140028922", false},
		{"(31) 40028922", "3140028922", false},
		{"(31) 98598-8922", "31985988922", false},
		{"(31) 985988922", "31985988922", false},
		{"123abc4567", "", true},
	}

	for _, test := range tests {
		phone := phone2.NewPhone(test.input)
		result, err := phone.AsUnmaskedString()

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input %s, but got nil", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
			} else if result != test.expected {
				t.Errorf("For input %s, expected %s but got %s", test.input, test.expected, result)
			}
		}
	}
}
