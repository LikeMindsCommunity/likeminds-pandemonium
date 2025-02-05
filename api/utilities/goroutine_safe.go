package utilities

import (
	"runtime/debug"

	"log"
)

// SafeGo starts a new goroutine with panic recovery
func SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Log the error and stack trace
				log.Println("Recovered from panic in goroutine: ", r)
				debug.PrintStack()
			}
		}()
		fn()
	}()
}
