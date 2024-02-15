package tests

import (
	"testing"

	"github.com/gtkmk/finder_api/core/domain/cpf"
)

func TestNewCpf(t *testing.T) {
	number := "12345678901"
	c := cpf.NewCpf(number)
	if c.Number != number {
		t.Errorf("Expected number %s, but got %s", number, c.Number)
	}
}

func TestCPF_AsString(t *testing.T) {
	number := "12345678901"
	c := cpf.NewCpf(number)
	if c.AsString() != number {
		t.Errorf("Expected AsString to return %s, but got %s", number, c.AsString())
	}
}

func TestCPF_AsUnmaskedString(t *testing.T) {
	number := "211.690.470-68"
	expected := "21169047068"
	c := cpf.NewCpf(number)
	unmasked, err := c.AsUnmaskedString()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if unmasked != expected {
		t.Errorf("Expected AsUnmaskedString to return %s, but got %s", expected, unmasked)
	}

	invalidNumber := "1234567810"
	c = cpf.NewCpf(invalidNumber)
	_, err = c.AsUnmaskedString()
	if err == nil {
		t.Error("Expected error for an invalid CPF, but got nil")
	}
}

func TestCPF_Validate(t *testing.T) {
	validNumber := "512.279.280-17"
	c := cpf.NewCpf(validNumber)
	err := c.Validate()
	if err != nil {
		t.Errorf("Unexpected error for a valid CPF: %v", err)
	}

	invalidNumber := "1234567810"
	c = cpf.NewCpf(invalidNumber)
	err = c.Validate()
	if err == nil {
		t.Error("Expected error for an invalid CPF, but got nil")
	}
}

func TestIsCpf(t *testing.T) {
	validNumber := "777.188.690-67"
	isValid := cpf.IsCpf(validNumber)
	if !isValid {
		t.Errorf("Expected IsCpf to return true for a valid CPF, but got false")
	}

	invalidNumber := "1234567810"
	isValid = cpf.IsCpf(invalidNumber)

	if isValid {
		t.Errorf("Expected IsCpf to return false for an invalid CPF, but got true")
	}
}
