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
