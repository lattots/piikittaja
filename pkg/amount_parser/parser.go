package amountparser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidAmount = errors.New("invalid amount format")
	ErrSubCentAmount = errors.New("sub-cent amounts (more than 2 decimals) are not allowed")
)

func ParseToCents(amount string) (int, error) {
	clean := strings.TrimSpace(amount)
	clean = strings.ReplaceAll(clean, "€", "")  // Remove € sign
	clean = strings.ReplaceAll(clean, ",", ".") // Convert , to . for English style flaot parsing
	clean = strings.TrimSpace(clean)

	if clean == "" {
		return 0, ErrInvalidAmount
	}

	// This will match 1st the sign, 2nd the whole part and 3rd the fraction part
	re := regexp.MustCompile(`^([-+]?)(\d*)?\.?(\d*)$`)
	matches := re.FindStringSubmatch(clean)
	if matches == nil || clean == "." || clean == "-" || clean == "+" {
		return 0, ErrInvalidAmount
	}

	sign := matches[1]
	wholeStr := matches[2]
	fracStr := matches[3]

	// If fraction part is more than 2 digits long, amount contains sub-cent value
	if len(fracStr) > 2 {
		return 0, ErrSubCentAmount
	}

	// If whole part does not exist (amount = .5), it is set to 0
	if wholeStr == "" {
		wholeStr = "0"
	}

	// Fraction is padded to two digits:
	// "5" -> "50"
	// "" -> "00"
	for len(fracStr) < 2 {
		fracStr += "0"
	}

	whole, err := strconv.Atoi(wholeStr)
	if err != nil {
		return 0, err
	}
	frac, err := strconv.Atoi(fracStr)
	if err != nil {
		return 0, err
	}

	// Amount parts are combined
	cents := (whole * 100) + frac

	if sign == "-" {
		cents = -cents
	}

	return cents, nil
}

func String(amount int) string {
	euros := amount / 100
	cents := amount % 100

	if cents == 0 {
		return fmt.Sprintf("%d €", euros)
	}

	absCents := cents
	if absCents < 0 {
		absCents = -absCents
	}

	return fmt.Sprintf("%d,%02d €", euros, absCents)
}

// GetCents returns the fractional/cent part of an amount.
// For example:
// 167 -> 67
// -8 -> -8
func GetCents(amount int) int {
	return amount % 100
}

// GetEuros returns the whole/euro part of an amount.
// For example:
// 167 -> 1
// -8 -> 0
func GetEuros(amount int) int {
	return amount / 100
}
