package pcw

// Field types

type PacketType uint8

const (
	PacketTypeKeyUp              PacketType = 0x00
	PacketTypeKeyDown            PacketType = 0x01
	PacketTypeElement            PacketType = 0x02
	PacketTypeCharacters         PacketType = 0x03
	PacketTypeWinKeyer           PacketType = 0x04
	PacketTypePing               PacketType = 0x05
	PacketTypePong               PacketType = 0x06
	PacketTypeMissed             PacketType = 0x07
	PacketTypeDropped            PacketType = 0x08
	PacketTypeSetJitterBuffer    PacketType = 0x09
	PacketTypeApplicationControl PacketType = 0x0A
)

// Interfaces

type Packet interface {
	Bytes() ([]byte, error)
}

// Functions

type DecodePacket func(header Header, payload []byte) (Packet, error)
