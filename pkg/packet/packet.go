package packet

import "time"

type Packet struct {
	Id       int
	Header   string
	ClientId int
	Expires  time.Time
}

func New(id int, header string, clientId int, expires time.Time) Packet {
	return Packet{
		Id:       id,
		Header:   header,
		ClientId: clientId,
		Expires:  expires,
	}
}
