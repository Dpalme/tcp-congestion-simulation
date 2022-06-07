package simulation

import (
	"math/rand"
	"tcp-congestion/pkg/client"
	congestionalgos "tcp-congestion/pkg/congestionAlgos"
	"tcp-congestion/pkg/logger"
	"tcp-congestion/pkg/server"
	"tcp-congestion/pkg/utils"
)

func RunFromCLI(congestionAlgo string, clientCount int, logger *logger.Logger) {
	congestionAlgorithm, err := congestionalgos.GetByName(congestionAlgo)
	if err != nil {
		logger.Critical(err.Error())
		return
	}
	s := server.New(congestionAlgorithm, *logger)

	for i := 0; i < clientCount; i++ {
		c := client.New(4 + rand.Intn(4))
		utils.TwoWayShake(c, s)
		go c.Run()
	}

	for {
		continue
	}
}
