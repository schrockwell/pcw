package pcw

import (
	"encoding/binary"
)

type KeyUp struct {
	Header    Header
	Timestamp uint32
}

func (p *KeyUp) Bytes() ([]byte, error) {
	err := validateEncode(p.Header, PacketTypeKeyUp)
	if err != nil {
		return nil, err
	}

	pb := make([]byte, 4)
	binary.BigEndian.PutUint32(pb[0:4], uint32(p.Timestamp))

	return buildPacketBuffer(p.Header, pb), nil
}

func DecodeKeyUp(h Header, pb []byte) (Packet, error) {
	err := validateDecode(h, pb, PacketTypeKeyUp, 4)
	if err != nil {
		return nil, err
	}

	return &KeyUp{
		Header:    h,
		Timestamp: binary.BigEndian.Uint32(pb[0:4]),
	}, nil
}
