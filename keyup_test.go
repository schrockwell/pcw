package pcw_test

import (
	"pcw"
	"pcw/test"
	"testing"
)

func TestBytes(t *testing.T) {
	pkt := &pcw.KeyUp{
		Header: pcw.Header{
			Type:    pcw.PacketTypeKeyUp,
			SeqNum:  1,
			Channel: 2,
		},
		Timestamp: 0x12345678,
	}

	hb, pb := test.TestPacketBytes(t, pkt)
	test.TestHeaderBytes(t, hb, pkt.Header)
	test.TestPayloadBytes(t, pb, []byte{0x12, 0x34, 0x56, 0x78})
}

func TestDecode(t *testing.T) {
	b := []byte{0, 1, 2, 0x12, 0x34, 0x56, 0x78}
	p := test.TestPacketDecode(t, b, pcw.DecodeKeyUp).(*pcw.KeyUp)

	test.TestHeader(t, p.Header, pcw.Header{
		Type:    pcw.PacketTypeKeyUp,
		SeqNum:  1,
		Channel: 2,
	})

	if p.Timestamp != 0x12345678 {
		t.Fatalf("expected timestamp %d, got %d", 0x12345678, p.Timestamp)
	}
}
