package TCP

import (
	congestioncontroller "tcp-congestion/pkg/congestionController"
	connection "tcp-congestion/pkg/connection"
)

var (
	SstMax      = 12.0
	SstInitial  float64
	CwndInitial float64
)

func New() (cc *congestioncontroller.CongestionControl) {
	return congestioncontroller.New("TCP", IncreaseProcedure, LossProcedure, TimeoutProcedure)
}

func IncreaseProcedure(conn *connection.Connection) {
	if conn.CWND < SstMax {
		conn.CWND++
		return
	}
	conn.CWND += 1 / (conn.CWND / (SstMax / 2))
}

func LossProcedure(conn *connection.Connection) {
	SstMax = conn.CWND
	conn.CWND /= 2
}

func TimeoutProcedure(conn *connection.Connection) {
	SstMax = conn.CWND
	conn.CWND = CwndInitial
}
