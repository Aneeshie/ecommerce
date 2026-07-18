package inventory

import "errors"

var (
	ErrNoProductFound = errors.New("No product found to get fetch inventory for.")
)
