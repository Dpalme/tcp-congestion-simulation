package client

import (
	"math/rand"
	"tcp-congestion/pkg/connection"
	"tcp-congestion/pkg/packet"
	"time"
)

type Client struct {
	IP    int
	delay time.Duration
	conn  *connection.Connection
}

func New(speed int) *Client {
	return &Client{
		IP:    rand.Intn(255 * 255 * 255 * 255),
		delay: time.Duration(rand.Float64() * float64(time.Second/time.Duration(speed))),
	}
}

func (c *Client) Run() {
	for p := range c.conn.Client.Read() {
		ackPacket := packet.New(p.Id, "ACK", p.ClientId, time.Now().Add(time.Second))
		c.conn.SendToServer(ackPacket)
		<-time.After(c.delay)
	}
}

func (c *Client) Connect(conn *connection.Connection) {
	c.conn = conn
}
