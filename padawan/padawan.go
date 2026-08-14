package padawan

import (
	"errors"
	"fmt"
)

var (
	DivisionByZeroError = errors.New("división by zero not allowed")
)

func Divide(dividend, divisor int) (int, error) {
	if divisor == 0 {
		return 0, fmt.Errorf("%w: [%d / %d]", DivisionByZeroError, dividend, divisor)
	}

	return dividend / divisor, nil
}
