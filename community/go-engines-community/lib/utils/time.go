package utils

import "time"

// DateOf returns time without hour, minute, and second.
func DateOf(v time.Time) time.Time {
	y, m, d := v.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, v.Location())
}

// MinTime returns minimal time between arguments.
func MinTime(left, right time.Time) time.Time {
	if left.Before(right) {
		return left
	}

	return right
}

// MaxTime returns max time between arguments.
func MaxTime(left, right time.Time) time.Time {
	if left.After(right) {
		return left
	}

	return right
}
