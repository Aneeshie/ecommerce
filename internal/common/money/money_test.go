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
		Name          string
		Left          int64
		Right         int64
		Expected      int64
		ExpectedError error
	}{
		{Name: "Positive + Positive", Left: 100, Right: 100, Expected: 200, ExpectedError: nil},
		{Name: "Positive + Zero", Left: 100, Right: 0, Expected: 100, ExpectedError: nil},
		{Name: "Zero + Zero", Left: 0, Right: 0, Expected: 0, ExpectedError: nil},
		{Name: "Overflow", Left: math.MaxInt64, Right: 1, Expected: 0, ExpectedError: ErrOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			left, _ := New(tt.Left)
			right, _ := New(tt.Right)

			result, err := left.Add(right)
			if err != tt.ExpectedError {
				t.Fatalf("expected error %v, got %v", tt.ExpectedError, err)
			}

			if result.Amount() != tt.Expected {
				t.Errorf("expected amount %d, got %d", tt.Expected, result.Amount())
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		Name          string
		Left          int64
		Right         int64
		Expected      int64
		ExpectedError error
	}{
		{Name: "Positive - Positive", Left: 200, Right: 100, Expected: 100, ExpectedError: nil},
		{Name: "Positive - Zero", Left: 100, Right: 0, Expected: 100, ExpectedError: nil},
		{Name: "Zero - Zero", Left: 0, Right: 0, Expected: 0, ExpectedError: nil},
		{Name: "Negative Result", Left: 100, Right: 200, Expected: 0, ExpectedError: ErrNegativeAmount},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			left, _ := New(tt.Left)
			right, _ := New(tt.Right)

			result, err := left.Subtract(right)
			if err != tt.ExpectedError {
				t.Fatalf("expected error %v, got %v", tt.ExpectedError, err)
			}

			if result.Amount() != tt.Expected {
				t.Errorf("expected amount %d, got %d", tt.Expected, result.Amount())
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		Name          string
		Amount        int64
		Multiplier    int
		Expected      int64
		ExpectedError error
	}{
		{Name: "Positive * Positive", Amount: 100, Multiplier: 2, Expected: 200, ExpectedError: nil},
		{Name: "Positive * Zero", Amount: 100, Multiplier: 0, Expected: 0, ExpectedError: nil},
		{Name: "Zero * Positive", Amount: 0, Multiplier: 5, Expected: 0, ExpectedError: nil},
		{Name: "Negative Multiplier", Amount: 100, Multiplier: -1, Expected: 0, ExpectedError: ErrNegativeAmount},
		{Name: "Overflow", Amount: math.MaxInt64, Multiplier: 2, Expected: 0, ExpectedError: ErrOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			money, _ := New(tt.Amount)

			result, err := money.Multiply(tt.Multiplier)
			if err != tt.ExpectedError {
				t.Fatalf("expected error %v, got %v", tt.ExpectedError, err)
			}

			if result.Amount() != tt.Expected {
				t.Errorf("expected amount %d, got %d", tt.Expected, result.Amount())
			}
		})
	}
}
