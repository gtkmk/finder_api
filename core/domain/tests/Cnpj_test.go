package tests

import (
	"testing"

	"github.com/gtkmk/finder_api/core/domain/cnpj"
)

func TestValidCNPJ(t *testing.T) {
	cnpjStr := "33.377.199/0001-20"
	cnpjObj := cnpj.NewCnpj(cnpjStr)
	err := cnpjObj.Validate()
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}
}

func TestInvalidCNPJ(t *testing.T) {
	cnpjStr := "00000000000000"
	cnpjObj := cnpj.NewCnpj(cnpjStr)
	err := cnpjObj.Validate()
	if err == nil {
		t.Error("Expected an error for invalid CNPJ, but got none")
	}

	cnpjStr = "1234567890"
	cnpjObj = cnpj.NewCnpj(cnpjStr)
	err = cnpjObj.Validate()
	if err == nil {
		t.Error("Expected an error for invalid CNPJ, but got none")
	}
}

func TestSanitize(t *testing.T) {
	input := "12.345.678/0001-90"
	expectedOutput := "12345678000190"
	cnpjObj := cnpj.NewCnpj(input)
	cnpjObj.Sanitize()
	if cnpjObj.Number != expectedOutput {
		t.Errorf("Expected %s, but got %s", expectedOutput, cnpjObj.Number)
	}
}
