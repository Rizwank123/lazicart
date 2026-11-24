package utils

import "time"

func ExponentialBackoff(attempt int) time.Duration {
	base := time.Second
	max := 30 * time.Second

	delay := base * time.Duration(1<<attempt)
	if delay > max {
		delay = max
	}
	return delay
}
