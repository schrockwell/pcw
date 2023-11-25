package pcw

func Decode(b []byte) (Packet, error) {
	h, pb, err := ParseHeader(b)
	if err != nil {
		return nil, err
	}

	switch h.Type {
	case PacketTypeKeyUp:
		return DecodeKeyUp(h, pb)
	default:
		return nil, ErrInvalidPacketType
	}
}

func buildPacketBuffer(h Header, pb []byte) []byte {
	hb := h.Bytes()

	buf := make([]byte, len(hb)+len(pb))
	copy(buf, hb)
	copy(buf[len(hb):], pb)

	return buf
}

func validateEncode(h Header, pt PacketType) error {
	if h.Type != pt {
		return ErrInvalidPacketType
	}

	return nil
}

func validateDecode(h Header, pb []byte, pt PacketType, l int) error {
	if h.Type != pt {
		return ErrInvalidPacketType
	}
	if len(pb) != l {
		return ErrInvalidPacket
	}

	return nil
}
