package money

import "errors"

type Money struct {
	paise int64
}

var ErrNegativeAmount = errors.New("money cannot be negative")

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

func (m Money) Add(other Money) Money {
	return Money{
		paise: m.paise + other.paise,
	}
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
