package tcpchannel

import (
	"tcp-congestion/pkg/packet"
)

type TCPChannel struct {
	ch          chan packet.Packet
	closed      bool
	PacketCount int
}

func New() *TCPChannel {
	return &TCPChannel{
		ch:          make(chan packet.Packet),
		closed:      false,
		PacketCount: 0,
	}
}

func (tc *TCPChannel) Close() {
	tc.closed = true
	close(tc.ch)
}

func (tc *TCPChannel) IsRunning() bool {
	return !tc.closed
}

func (tc *TCPChannel) SendPacket(packet packet.Packet) bool {
	if tc.closed {
		return false
	}
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	tc.ch <- packet
	return true
}

func (tc *TCPChannel) Read() <-chan packet.Packet {
	return tc.ch
}
