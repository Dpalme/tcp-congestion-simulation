package congestiontest

import (
	"fmt"
	"math/rand"
	"tcp-congestion/pkg/client"
	congestioncontroller "tcp-congestion/pkg/congestionController"
	"tcp-congestion/pkg/logger"
	"tcp-congestion/pkg/server"
	"tcp-congestion/pkg/utils"
	"time"
)

type CongestionTest struct {
	cc              *congestioncontroller.CongestionControl
	clientCount     int64
	testDuration    float64
	PacketsSent     int
	PacketsReceived int
	logger          *logger.Logger
	currentRun      int
}

func New(cc *congestioncontroller.CongestionControl, clientCount int64, testDuration float64, logger *logger.Logger) *CongestionTest {
	return &CongestionTest{
		cc:           cc,
		clientCount:  clientCount,
		testDuration: testDuration,
		logger:       logger,
		currentRun:   0,
	}
}

func (test *CongestionTest) Run() {
	test.logger.SetPrefix(fmt.Sprintf("Run %d: ", test.currentRun))
	test.logger.Default("Starting %s run with %d clients for %.2f seconds", test.cc.Name, test.clientCount, test.testDuration)
	s := server.New(test.cc, *test.logger)

	for i := int64(0); i < test.clientCount; i++ {
		c := client.New(4 + rand.Intn(4))
		utils.TwoWayShake(c, s)
		go c.Run()
	}

	defer func() {
		if r := recover(); r != nil {
			test.logger.Default("Recovered from deadlock")
		}
	}()

	<-time.After(time.Second * time.Duration(test.testDuration))
	s.Terminate()
	test.logger.SetPrefix("")
	sent, received, _, _ := s.GetStats()
	test.PacketsSent += sent
	test.PacketsReceived += received

	test.PrintRunStats(s)
}

func (test *CongestionTest) CalculateRunStats(s *server.Server) (sent int, received int, sentPerClientPerSecond float64, receivedPerClientPerSecond float64, inTransit int, avgCWND float64) {
	sent, received, inTransit, avgCWND = s.GetStats()
	sentPerClientPerSecond = float64(sent) / float64(test.clientCount) / float64(test.testDuration)
	receivedPerClientPerSecond = float64(received) / float64(test.clientCount) / float64(test.testDuration)

	return sent, received, sentPerClientPerSecond, receivedPerClientPerSecond, inTransit, avgCWND
}

func (test *CongestionTest) PrintRunStats(s *server.Server) {
	sent, received, sentPerClientPerSecond, receivedPerClientPerSecond, inTransit, avgCWND := test.CalculateRunStats(s)
	test.logger.Informational(fmt.Sprintf("%10d packets sent\n", sent) +
		fmt.Sprintf("%10.2f packets sent per second\n", sentPerClientPerSecond) +
		fmt.Sprintf("%10d packets received\n", received) +
		fmt.Sprintf("%10.2f packets received per second\n", receivedPerClientPerSecond) +
		fmt.Sprintf("%10d packets left in transit\n", inTransit) +
		fmt.Sprintf("%10.2f avg CWND", avgCWND))

}

func (test *CongestionTest) RunNTimes(repeat int) {
	for i := 0; i < repeat; i++ {
		test.Run()
	}
	test.printOutput(repeat)
}

func (test *CongestionTest) printOutput(runs int) {
	avgSent := float64(test.PacketsSent) / float64(runs)
	avgreceived := float64(test.PacketsReceived) / float64(runs)

	test.logger.Default("\nFinished %d runs", runs)
	test.logger.Default("%10d packets sent", test.PacketsSent)
	test.logger.Default("%10.2f avg/run", avgSent)
	test.logger.Default("%10.4f/client", avgSent/float64(test.clientCount))
	test.logger.Default("%10d packets received", test.PacketsReceived)
	test.logger.Default("%10.2f avg/run", avgreceived)
	test.logger.Default("%10.4f/client)", avgreceived/float64(test.clientCount))
	test.logger.Default("That's %.2f%% fulfilled", (avgreceived/avgSent)*100)
}
