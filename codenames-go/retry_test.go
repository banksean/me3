package main

import (
	"fmt"
	"testing"
)

func TestRetry_Success(t *testing.T) {
	err := retry(1, func() error {
		return nil
	})
	if err != nil {
		t.Errorf("expected nil error but got %v", err)
	}
}

func TestRetry_SuccessRetries(t *testing.T) {
	count := 2
	err := retry(2, func() error {
		count--
		if count > 0 {
			return fmt.Errorf("won't succeed yet, counter isn't 0: %d", count)
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected nil error but got: %v", err)
	}
}

func TestRetry_FailRetries(t *testing.T) {
	count := 3
	err := retry(2, func() error {
		count--
		if count > 0 {
			return fmt.Errorf("won't succeed yet, counter isn't 0: %d", count)
		}
		return nil
	})
	if err == nil {
		t.Errorf("retry should have returned a non-nil error")
	}
}
