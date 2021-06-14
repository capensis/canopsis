package utils

import "log"

// FailOnError calls log.Panicf if err != nil
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %v", msg, err)
	}
}

// ExitOnError calls log.Fatalf if err != nil
func ExitOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
