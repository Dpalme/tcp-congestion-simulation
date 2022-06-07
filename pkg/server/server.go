package server

import (
	"fmt"
	"math/rand"
	"sync"
	congestioncontroller "tcp-congestion/pkg/congestionController"
	"tcp-congestion/pkg/connection"
	"tcp-congestion/pkg/logger"
	"tcp-congestion/pkg/packet"
	"time"
)

type Server struct {
	Cc_algo     *congestioncontroller.CongestionControl
	timeout     time.Duration
	Connections map[int]*connection.Connection
	cmu         *sync.Mutex
	log         logger.Logger
}

func New(cc *congestioncontroller.CongestionControl, logger logger.Logger) (s *Server) {
	return &Server{
		timeout:     time.Second * 2,
		Connections: make(map[int]*connection.Connection),
		cmu:         &sync.Mutex{},
		log:         logger,
		Cc_algo:     cc,
	}
}

func (s *Server) GetConnection(id int) *connection.Connection {
	s.cmu.Lock()
	defer s.cmu.Unlock()
	return s.Connections[id]
}

func (s *Server) SetConnection(id int, conn *connection.Connection) {
	s.cmu.Lock()
	defer s.cmu.Unlock()
	s.Connections[id] = conn
}

func (s *Server) Terminate() {
	s.log.Default("Finished Run")

	s.log.Default("Closing %d connections", len(s.Connections))
	for _, conn := range s.Connections {
		s.log.Informational("%10.2f CWND", conn.CWND)
		conn.Terminate()
	}
}

func (s *Server) AddConnection(conn *connection.Connection) {
	s.SetConnection(conn.Id, conn)
	s.log.Default("Connected to Client[%d]", conn.Id)

	s.infiniteData(conn)
}

func (s *Server) infiniteData(conn *connection.Connection) {
	go func(conn *connection.Connection) {
		for packet := range conn.Server.Read() {
			s.log.Debug("Got %s Packet[%d] from Client[%d]", packet.Header, packet.Id, packet.ClientId)
			s.handlePacket(packet)
		}
	}(conn)

	go func(conn *connection.Connection) {
		for conn.Client.IsRunning() {
			s.log.Default("Connection[%10d] with CWND %10.2f", conn.Id, conn.CWND)
			if float64(conn.InTransitLength()) < conn.CWND {
				newPacket := packet.New(rand.Int(), "DATA", conn.Id, time.Now().Add(s.timeout))

				s.log.Debug("Sending Packet[%d] to Client[%d]", newPacket.Id, newPacket.ClientId)

				if !conn.SendToClient(newPacket) {
					return
				}

				go s.checkLoss(conn, newPacket)
			}
		}
	}(conn)
}

func (s *Server) handlePacket(packet packet.Packet) {
	conn := s.GetConnection(packet.ClientId)

	if conn.IsInTransit(packet.Id) {
		conn.PacketArrived(packet.Id)
		s.Cc_algo.IncreaseProcedure(conn)
	}
}

func (s *Server) LogStats() {
	sent, received, inTransit, avgCWND := s.GetStats()
	s.log.Informational(fmt.Sprintf("%10d packets sent\n", sent) +
		fmt.Sprintf("%10d packets received\n", received) +
		fmt.Sprintf("%10d packets left in transit\n", inTransit) +
		fmt.Sprintf("%5.2f avg CWND", avgCWND))
}

func (s *Server) GetStats() (sent int, received int, inTransit int, avgCWND float64) {
	for _, conn := range s.Connections {
		inTransit += conn.InTransitLength()
		avgCWND += conn.CWND
		sent += conn.Server.PacketCount
		received += conn.Client.PacketCount
	}

	avgCWND /= float64(len(s.Connections))
	return sent, received, inTransit, avgCWND
}

func (s *Server) checkLoss(conn *connection.Connection, data packet.Packet) {
	<-time.After(s.timeout)
	if !conn.IsInTransit(data.Id) {
		return
	}

	s.log.Debug("Packet[%d] for Client[%d] lost", data.Id, data.ClientId)

	conn.IsInTransit(data.Id)
	s.Cc_algo.LossProcedure(conn)
	go s.checkTimeout(conn, data)
}

func (s *Server) checkTimeout(conn *connection.Connection, data packet.Packet) {
	<-time.After(s.timeout)
	if !conn.IsInTransit(data.Id) {
		return
	}

	s.log.Debug("Packet[%d] for Client[%d] timed out", data.Id, data.ClientId)

	conn.PacketArrived(data.Id)
	s.Cc_algo.TimeoutProcedure(conn)
}
