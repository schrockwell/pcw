package pcw

type Header struct {
	Type    PacketType
	SeqNum  uint8
	Channel uint8
}

const HeaderLength = 3

func (h *Header) Bytes() []byte {
	return []byte{byte(h.Type), byte(h.SeqNum), byte(h.Channel)}
}

func ParseHeader(b []byte) (Header, []byte, error) {
	if len(b) < HeaderLength {
		return Header{}, nil, ErrInvalidHeader
	}

	hb, pb := SplitPacketBytes(b)

	return Header{
		Type:    PacketType(hb[0]),
		SeqNum:  hb[1],
		Channel: hb[2],
	}, pb, nil
}

func SplitPacketBytes(b []byte) ([]byte, []byte) {
	return b[:HeaderLength], b[HeaderLength:]
}
