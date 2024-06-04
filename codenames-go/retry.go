package main

import (
	"fmt"
	"log"
)

func retry(attempts int, f func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error:", err)
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
