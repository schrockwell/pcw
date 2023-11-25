package pcw

// Field types

type PacketType uint8
type Microseconds uint32 // microseconds

// Interfaces

type Packet interface {
	Bytes() ([]byte, error)
}

// Functions

type DecodePacket func(header Header, payload []byte) (Packet, error)
