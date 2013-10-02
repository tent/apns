package apns

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Notification struct {
	Token    []byte
	Payload  []byte
	ID       uint32
	Expiry   time.Time
	Priority byte
}

func (n *Notification) length() int {
	return 1 + // command
		4 + // frame length
		3 + // device token item header
		32 + // device token
		3 + // payload item header
		len(n.Payload) + // payload
		3 + // identifier item header
		4 + // identifier
		3 + // expiration item header
		4 + // expiration
		3 + // priority item header
		1 // priority
}

func (n *Notification) Bytes() []byte {
	length := n.length()
	buf := bytes.NewBuffer(make([]byte, 0, length))

	// Command
	buf.WriteByte(2)
	// Frame Length
	binary.Write(buf, binary.BigEndian, uint32(length-5))

	// Device Token Item
	buf.WriteByte(1)
	// Device Token Length
	binary.Write(buf, binary.BigEndian, uint16(32))
	// Device Token
	buf.Write(n.Token)

	// Payload Item
	buf.WriteByte(2)
	// Payload Length
	binary.Write(buf, binary.BigEndian, uint16(len(n.Payload)))
	// Payload
	buf.Write(n.Payload)

	// Identifier Item
	buf.WriteByte(3)
	// Identifier Length
	binary.Write(buf, binary.BigEndian, uint16(4))
	// Identifier
	binary.Write(buf, binary.BigEndian, n.ID)

	// Expiration Item
	buf.WriteByte(4)
	// Expiration Length
	binary.Write(buf, binary.BigEndian, uint16(4))
	// Expiration
	if n.Expiry.IsZero() {
		binary.Write(buf, binary.BigEndian, uint32(0))
	} else {
		binary.Write(buf, binary.BigEndian, uint32(n.Expiry.Unix()))
	}

	if n.Priority == 0 {
		n.Priority = 10
	}
	// Priority Item
	buf.WriteByte(5)
	// Priority Length
	binary.Write(buf, binary.BigEndian, uint16(1))
	// Priority
	buf.WriteByte(n.Priority)

	return buf.Bytes()
}
