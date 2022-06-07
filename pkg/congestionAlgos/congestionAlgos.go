package congestionalgos

import (
	"fmt"
	"strings"
	"tcp-congestion/pkg/congestionAlgos/BIC"
	"tcp-congestion/pkg/congestionAlgos/CTCP"
	"tcp-congestion/pkg/congestionAlgos/NONE"
	"tcp-congestion/pkg/congestionAlgos/TCP"
	congestioncontroller "tcp-congestion/pkg/congestionController"
)

var (
	RegisteredAlgorithms = []string{"TCP", "BIC", "CTCP", "NONE"}
)

func GetByName(name string) (*congestioncontroller.CongestionControl, error) {
	switch strings.ToUpper(name) {
	case "TCP":
		return TCP.New(), nil
	case "BIC":
		return BIC.New(), nil
	case "CTCP":
		return CTCP.New(), nil
	case "NONE":
		return NONE.New(), nil
	}
	return nil, fmt.Errorf("unknown congestion algorithm: %s", name)
}
