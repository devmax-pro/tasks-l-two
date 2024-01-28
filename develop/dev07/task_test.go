package main

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-or(
		sig(1*time.Second),
		sig(5*time.Second),
		sig(10*time.Second),
		sig(1*time.Minute),
		sig(1*time.Hour),
	)

	duration := time.Since(start)

	if duration.Seconds() > 1.1 {
		t.Errorf("Expected <= 1.1, got %v", duration.Seconds())
	}
}
