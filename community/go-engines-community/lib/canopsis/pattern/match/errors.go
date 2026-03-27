package match

import "fmt"

// ConditionError represents a validation error for a specific condition within a pattern group.
type ConditionError struct {
	GroupIdx int
	CondIdx  int
	Err      error
}

func (e ConditionError) Error() string {
	return fmt.Sprintf("invalid condition [%d][%d]: %s", e.GroupIdx, e.CondIdx, e.Err.Error())
}

func (e ConditionError) Unwrap() error {
	return e.Err
}
