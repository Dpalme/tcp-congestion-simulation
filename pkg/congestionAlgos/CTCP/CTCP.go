package CTCP

import (
	tcp "tcp-congestion/pkg/congestionAlgos/TCP"
	congestioncontroller "tcp-congestion/pkg/congestionController"
	"tcp-congestion/pkg/connection"
	"time"
)

var (
	aC       = 0.0
	beta_C   = 0.5
	cdc      = false
	gamma_C  = 30.0
	dc       = 0.0
	dwnd     = 0.0
	ec       = 0.0
	zeta_c   = 0.1
	lwc      = 41.0
	minRTTC  = 0.0
	srttc    = time.Second
	cwndint  = 1.0
	int_max  = 65535.0
	running  = false
	lastTrip = time.Now()
)

func New() *congestioncontroller.CongestionControl {
	return congestioncontroller.New("CTCP", IncreaseProcedure, TimeoutProcedure, LossProcedure)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func periodicUpdate(conn *connection.Connection) {
	ec = conn.CWND / minRTTC
	aC = conn.CWND / float64(srttc)
	dc = (ec - aC) * minRTTC
	if cdc {
		dwnd = min(0, conn.CWND*(1-beta_C)-(conn.CWND/2))
		cdc = false
	}
	if !cdc || dc < gamma_C {
		dwnd += min(0, gamma_C*conn.CWND-1)
	} else {
		dwnd += min(0, dwnd*(zeta_c-dc))
	}
	conn.CWND = min(int_max, conn.CWND+dwnd)
}

func inifiteUpdate(conn *connection.Connection) {
	for {
		<-time.After(time.Duration(srttc))
		periodicUpdate(conn)
	}
}

func updatesrttc() {
	rtt := time.Since(lastTrip)
	srttc = (rtt + srttc) / 2
	lastTrip = time.Now()
}

func IncreaseProcedure(conn *connection.Connection) {
	if conn.CWND < lwc {
		tcp.IncreaseProcedure(conn)
		return
	}
	conn.CWND += (1 / (conn.CWND + float64(dwnd)))
	conn.CWND += float64(dwnd)
	if !running {
		go inifiteUpdate(conn)
	}
	updatesrttc()
}

func LossProcedure(conn *connection.Connection) {
	tcp.SstMax = conn.CWND
	conn.CWND = (conn.CWND / 2.0) + float64(dwnd)
	cdc = true
	updatesrttc()
}

func TimeoutProcedure(conn *connection.Connection) {
	tcp.SstMax = max(conn.CWND/2, cwndint)
	conn.CWND = cwndint
	dwnd = 0
	cdc = true
	updatesrttc()
}
