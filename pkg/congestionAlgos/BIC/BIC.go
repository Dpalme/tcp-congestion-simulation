package BIC

import (
	tcp "tcp-congestion/pkg/congestionAlgos/TCP"
	congestioncontroller "tcp-congestion/pkg/congestionController"
	"tcp-congestion/pkg/connection"
)

var (
	BetaB     = 0.8
	Bb        = 4.0
	DeltaB    float64
	LWB       = 14.0
	Min       float64
	Max       float64
	Prev      float64
	sigma     = 20.0
	SMAX      = 32.0
	ss        = true
	ss_target float64
	ss_window float64
	target    float64
)

func New() *congestioncontroller.CongestionControl {
	return congestioncontroller.New("BIC", IncreaseProcedure, TimeoutProcedure, LossProcedure)
}

func f_alpha(delta float64, CWND float64, target float64) float64 {
	if (delta < 1 && CWND < target) || target <= CWND && CWND < (target+Bb) {
		return Bb / sigma
	}
	if (1 < delta && delta < SMAX) && CWND < target {
		return delta
	}
	if Bb < (CWND-target) && (CWND-target) < (SMAX*(Bb-1)) {
		return target / (Bb - 1)
	}
	return SMAX

}

func IncreaseProcedure(conn *connection.Connection) {
	if conn.CWND < LWB {
		tcp.IncreaseProcedure(conn)
		return
	}
	if ss {
		DeltaB = (target - conn.CWND) / Bb
		conn.CWND += f_alpha(DeltaB, conn.CWND, target) / conn.CWND
		if conn.CWND < Max {
			Min = conn.CWND
			target = (Max + Min) / 2
		} else {
			ss = true
			Max = conn.CWND * 2
			ss_window = 1
			ss_target = conn.CWND + 1
		}
	} else {
		conn.CWND += ss_window / conn.CWND
		if conn.CWND >= ss_target {
			ss_window *= Bb / (Bb - 1)
			ss_target = conn.CWND + ss_window
		}
		if ss_window >= Max {
			ss_window += 1
		}
	}
}

func LossProcedure(conn *connection.Connection) {
	if conn.CWND < Prev {
		Max = ((1.0 + BetaB) / 2.0) * conn.CWND
	} else {
		Max = conn.CWND
	}

	Prev = Max
	target = 0
	ss = false
	conn.CWND *= BetaB
	ss_target = conn.CWND
}

func TimeoutProcedure(conn *connection.Connection) {
	if conn.CWND < Prev {
		Max = ((1 + BetaB) / 2) * conn.CWND
	} else {
		Max = conn.CWND
	}

	Prev = Max
	target = 0
	ss = false
	conn.CWND *= BetaB
	ss_target = conn.CWND
}
