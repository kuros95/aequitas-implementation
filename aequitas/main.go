package aeuqitas

import "time"

func TimeCheck(goal time.Duration, actual time.Duration) bool {
	if actual > goal {
		return false
	} else {
		return true
	}
}

// 1. intercept traffic
// 2. start time measuring
// 3. on rpc complete stop measuring time
// 4. lower priority if time exceeded
