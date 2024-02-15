package cpf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	CPF_LENGTH                       = 11
	CpfIsInvalidErrorConst           = "O CPF é inválido."
	InvalidCpfLengthErrorConst       = "O CPF não possui o tamanho correto."
	CpfShouldContainOnlyNumbersConts = "O CPF deve conter apenas números."
)

type CPF struct {
	Number string
}

func NewCpf(number string) *CPF {
	return &CPF{
		number,
	}
}

func (cpf *CPF) AsString() string {
	return cpf.Number
}

func (cpf *CPF) AsUnmaskedString() (string, error) {
	if err := cpf.Validate(); err != nil {
		return "", err
	}

	return cpf.Number, nil
}

func (cpf *CPF) Validate() error {
	cpf.Unmask()

	cpf.Number = strings.Trim(cpf.Number, " ")
	cpf.Number = fmt.Sprintf("%011s", cpf.Number)

	_, err := strconv.Atoi(cpf.Number)

	if err != nil {
		return fmt.Errorf(CpfShouldContainOnlyNumbersConts)
	}

	if err := cpf.isValid(); err != nil {
		return err
	}

	return nil
}

func (cpf *CPF) Unmask() {
	re := regexp.MustCompile(`[.-]`)
	cpf.Number = re.ReplaceAllString(cpf.Number, "")
}

func IsCpf(number string) bool {
	if err := NewCpf(number).Validate(); err != nil {
		return false
	}

	return true
}

func (cpf *CPF) isValid() error {
	if err := cpf.checkCpfValidity(); err != nil {
		return err
	}

	return nil
}

func (cpf *CPF) checkCpfValidity() error {
	cpf.Number = regexp.MustCompile(`\D`).ReplaceAllString(cpf.Number, "")

	if len(cpf.Number) != CPF_LENGTH {
		return fmt.Errorf(InvalidCpfLengthErrorConst)
	}

	digits, err := cpf.parseCpfDigits()
	if err != nil {
		return err
	}

	mod1, mod2 := cpf.calculateChecksums(digits)

	if mod1 != digits[9] || mod2 != digits[10] {
		return fmt.Errorf(CpfIsInvalidErrorConst)
	}

	return nil
}

func (cpf *CPF) parseCpfDigits() ([]int, error) {
	digits := make([]int, CPF_LENGTH)
	for i, char := range cpf.Number {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return nil, fmt.Errorf(CpfIsInvalidErrorConst)
		}
		digits[i] = digit
	}
	return digits, nil
}

func (cpf *CPF) calculateChecksums(digits []int) (int, int) {
	sum1 := 0
	for i := 0; i < 9; i++ {
		sum1 += digits[i] * (10 - i)
	}
	mod1 := (sum1 * 10) % 11
	if mod1 == 10 {
		mod1 = 0
	}

	sum2 := 0
	for i := 0; i < 9; i++ {
		sum2 += digits[i] * (11 - i)
	}
	sum2 += mod1 * 2
	mod2 := (sum2 * 10) % 11
	if mod2 == 10 {
		mod2 = 0
	}

	return mod1, mod2
}

func (cpf *CPF) AddMask() string {
	re := regexp.MustCompile(`(\d{3})(\d{3})(\d{3})(\d{2})`)
	return re.ReplaceAllString(cpf.Number, "$1.$2.$3-$4")
}
