package utils

import (
	"fmt"
	"time"
)

func AequitasInit(latency_target int, target_pctl int) {
	for no_of_prio := range len(prios) {
		incr_window := latency_target * (100 / (100 - target_pctl))
		prios[no_of_prio].incr_window = incr_window
	}
}

func lowerPrio(goal time.Duration, elapsed time.Duration) bool {
	if elapsed > goal {
		fmt.Println("priority lowered")
		return true
	} else {
		return false
	}
}

func (r rpc) admit() bool {
	if r.prio.prio == "hi" {
		r.goal = 20 * time.Millisecond
	} else {
		r.goal = 15 * time.Millisecond
	}
	return lowerPrio(r.goal, r.elapsed)
}
