package NONE

import (
	congestioncontroller "tcp-congestion/pkg/congestionController"
	"tcp-congestion/pkg/connection"
)

var (
	int_max = 65535.0
)

func New() *congestioncontroller.CongestionControl {
	return congestioncontroller.New("None", IncreaseProcedure, TimeoutProcedure, LossProcedure)
}

func IncreaseProcedure(conn *connection.Connection) {
	conn.CWND = int_max
}

func LossProcedure(conn *connection.Connection) {

}

func TimeoutProcedure(conn *connection.Connection) {

}
