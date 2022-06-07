package connection

import (
	"sync"
	"tcp-congestion/pkg/packet"
	tcpchannel "tcp-congestion/pkg/tcpChannel"
)

type Connection struct {
	Id        int
	CWND      float64
	Server    *tcpchannel.TCPChannel
	Client    *tcpchannel.TCPChannel
	inTransit map[int]bool
	mu        *sync.Mutex
}

func New(IP int) *Connection {
	return &Connection{
		Id:        IP,
		CWND:      1,
		Server:    tcpchannel.New(),
		Client:    tcpchannel.New(),
		inTransit: make(map[int]bool),
		mu:        &sync.Mutex{},
	}
}

func (c *Connection) IsInTransit(packetId int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.inTransit[packetId]
	return ok && val
}

func (c *Connection) InTransitLength() int {
	return len(c.inTransit)
}

func (c *Connection) Terminate() {
	go c.Server.Close()
	go c.Client.Close()
}

func (c *Connection) SendToClient(packet packet.Packet) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Server.PacketCount++
	c.inTransit[packet.Id] = true
	return c.Client.SendPacket(packet)
}

func (c *Connection) SendToServer(packet packet.Packet) bool {
	return c.Server.SendPacket(packet)
}

func (c *Connection) PacketArrived(packetId int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists := c.inTransit[packetId]
	if exists {
		delete(c.inTransit, packetId)
		c.Client.PacketCount++
	}
}
