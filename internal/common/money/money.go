package money

import "errors"

type Money struct {
	paise int64
}

var ErrNegativeAmount = errors.New("money cannot be negative")
var ErrOverflow = errors.New("Amount too large")

func New(amount int64) (Money, error) {
	if amount < 0 {
		return Money{}, ErrNegativeAmount
	}

	return Money{paise: amount}, nil
}

func (m Money) Amount() int64 {
	return m.paise
}

func Zero() Money {
	return Money{}
}

func (m Money) Add(other Money) (Money, error) {
	sum := m.paise + other.paise
	overflow := ((m.paise ^ sum) & (other.paise ^ sum)) < 0

	if overflow {
		return Money{}, ErrOverflow
	}

	return Money{
		paise: m.paise + other.paise,
	}, nil
}

func (m Money) Subtract(other Money) Money {
	return Money{
		paise: m.paise - other.paise,
	}
}

func (m Money) Multiply(quantity int) Money {
	return Money{
		paise: m.paise * int64(quantity),
	}
}
