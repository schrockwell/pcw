package pcw

// Field types

type PacketType uint8
type Microseconds uint32 // microseconds

// Interfaces

type Encoder interface {
	Bytes() ([]byte, error)
}

type Packet interface {
	Encoder
}

// Functions

type DecodePacket func(header Header, payload []byte) (Packet, error)
