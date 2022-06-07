package congestioncontroller

import "tcp-congestion/pkg/connection"

type CongestionControl struct {
	IncreaseProcedure func(*connection.Connection)
	TimeoutProcedure  func(*connection.Connection)
	LossProcedure     func(*connection.Connection)
	Name              string
}

func New(name string, increaseProcedure func(*connection.Connection), timeoutProcedure func(*connection.Connection), lossProcedure func(*connection.Connection)) (cc *CongestionControl) {
	return &CongestionControl{
		IncreaseProcedure: increaseProcedure,
		TimeoutProcedure:  timeoutProcedure,
		LossProcedure:     lossProcedure,
		Name:              name,
	}
}
