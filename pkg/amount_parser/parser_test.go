// This is an auto generated test file for amountparser

package amountparser_test

import (
	"errors"
	"testing"

	amountparser "github.com/lattots/piikittaja/pkg/amount_parser"
)

func TestParseToCents(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr error
	}{
		// Basic success cases
		{"Standard amount", "1.25", 125, nil},
		{"Whole number", "10", 1000, nil},
		{"Zero decimal", "5.00", 500, nil},
		{"Single decimal (padding)", "1.1", 110, nil},
		{"Leading dot", ".50", 50, nil},
		{"Leading dot single digit", ".5", 50, nil},

		// Signs and Symbols
		{"Negative amount", "-1.25", -125, nil},
		{"Positive sign", "+1.25", 125, nil},
		{"With Euro symbol", "€1.25", 125, nil},
		{"With Euro and space", "€ 10.50", 1050, nil},

		// Localized separators
		{"Comma as decimal", "1,50", 150, nil},
		{"Comma single digit", "1,5", 150, nil},

		// Error cases: Invalid formats
		{"Empty string", "", 0, amountparser.ErrInvalidAmount},
		{"Just a dot", ".", 0, amountparser.ErrInvalidAmount},
		{"Just a sign", "-", 0, amountparser.ErrInvalidAmount},
		{"Alphabetic characters", "12abc", 0, amountparser.ErrInvalidAmount},
		{"Two dots", "1.2.3", 0, amountparser.ErrInvalidAmount},

		// Error cases: Precision
		{"Sub-cent amount (3 digits)", "1.001", 0, amountparser.ErrSubCentAmount},
		{"Sub-cent amount (4 digits)", "0.0001", 0, amountparser.ErrSubCentAmount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := amountparser.ParseToCents(tt.input)

			// Check for error expectations
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ParseToCents() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			// Check for unexpected errors
			if err != nil {
				t.Fatalf("ParseToCents() unexpected error: %v", err)
			}

			// Check for correct value
			if got != tt.want {
				t.Errorf("ParseToCents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		amount int
		str    string
	}{
		{100, "1 €"},
		{150, "1,50 €"},
		{101, "1,01 €"},
		{0, "0 €"},
		{-100, "-1 €"},
		{-101, "-1,01 €"},
	}

	for _, testCase := range testCases {
		str := amountparser.String(testCase.amount)
		if str != testCase.str {
			t.Errorf("want \"%s\", got \"%s\"\n", testCase.str, str)
		}
	}
}

func TestGetParts(t *testing.T) {
	testCases := []struct {
		original int
		euros    int
		cents    int
	}{
		{100, 1, 0},
		{-100, -1, 0},
		{150, 1, 50},
		{-150, -1, -50},
	}

	for _, testCase := range testCases {
		euros := amountparser.GetEuros(testCase.original)
		cents := amountparser.GetCents(testCase.original)
		if euros != testCase.euros {
			t.Errorf("want %d got %d", testCase.euros, euros)
		}
		if cents != testCase.cents {
			t.Errorf("want %d got %d", testCase.cents, cents)
		}
	}
}
