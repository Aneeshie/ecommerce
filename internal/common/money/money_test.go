package money

import (
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		Name           string
		Input          int64
		ExpectedAmount int64
		ExpectedError  error
	}{
		{Name: "Positive Amount", Input: 100, ExpectedAmount: 100, ExpectedError: nil},
		{Name: "Zero Amount", Input: 0, ExpectedAmount: 0, ExpectedError: nil},
		{Name: "Negative Amount", Input: -100, ExpectedAmount: 0, ExpectedError: ErrNegativeAmount},
		{Name: "Max Int", Input: math.MaxInt64, ExpectedAmount: math.MaxInt64, ExpectedError: nil},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			money, err := New(tt.Input)
			if err != tt.ExpectedError {
				t.Fatalf("expected error %v, got %v", err, tt.ExpectedError)
			}

			if money.Amount() != tt.ExpectedAmount {
				t.Errorf("expected amount %v, got %v", tt.ExpectedAmount, money.Amount())
			}
		})

	}
}
func TestAdd(t *testing.T) {
	tests := []struct {
		Name     string
		Left     int64
		Right    int64
		Expected int64
	}{
		{
			Name:     "Positive + Positive",
			Left:     100,
			Right:    100,
			Expected: 200,
		},
		{
			Name:     "Positive + Zero",
			Left:     100,
			Right:    0,
			Expected: 100,
		},
		{
			Name:     "Zero + Zero",
			Left:     0,
			Right:    0,
			Expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			left, err := New(tt.Left)
			if err != nil {
				t.Fatalf("failed to create left money: %v", err)
			}

			right, err := New(tt.Right)
			if err != nil {
				t.Fatalf("failed to create right money: %v", err)
			}

			result := left.Add(right)

			if result.Amount() != tt.Expected {
				t.Errorf("expected %d, got %d", tt.Expected, result.Amount())
			}
		})
	}
}
func TestAddOverflow(t *testing.T) {
	left, _ := New(math.MaxInt64)
	right, _ := New(1)

	result := left.Add(right)

	t.Logf("result = %d", result.Amount())
}
