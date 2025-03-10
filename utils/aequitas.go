package utils

import (
	"time"
)

func AequitasInit(latency_target int, target_pctl int) {
	for prio := range len(prios) {
		incr_window := latency_target * (100 / (100 - target_pctl))
		prios[prio].incr_window = time.Duration(incr_window) * time.Millisecond
	}
}

func (r rpc) admit(additive_increase, multiplicative_decrease, minimum_admission float64) {
	if r.prio.latency > (r.elapsed / time.Duration(r.size)) {
		if time.Since(time.Now())-time.Since(r.prio.t_last_increase) > r.incr_window {
			r.prio.p_admit = min(r.prio.p_admit+additive_increase, 1)
		}
		r.t_last_increase = time.Now()
	} else {
		//the second value is the minimum possible probability of admission for a given RPC,
		//since it is not desirable to have 0 probability, as it would starve the network
		r.prio.p_admit = max(r.prio.p_admit-(multiplicative_decrease*float64(r.size)), minimum_admission)
	}
}
