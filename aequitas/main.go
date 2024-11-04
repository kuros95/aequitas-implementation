package aequitas

import (
	"fmt"
	"time"
)

func LowerPrio(goal time.Duration, elapsed time.Duration) bool {
	if elapsed > goal {
		fmt.Println("priority lowered")
		return true
	} else {
		return false
	}
}

// 1. intercept traffic
// 2. start time measuring
// 3. on rpc complete stop measuring time
// 4. lower priority if time exceeded
