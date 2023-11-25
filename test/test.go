package test

import (
	"pcw"
	"reflect"
	"testing"
)

func TestPacketBytes(t *testing.T, p pcw.Packet) ([]byte, []byte) {
	b, err := p.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	return pcw.SplitPacketBytes(b)
}

func TestHeaderBytes(t *testing.T, buf []byte, h pcw.Header) {
	if len(buf) != 3 {
		t.Fatalf("expected 3 bytes, got %d", len(buf))
	}

	if buf[0] != byte(h.Type) {
		t.Fatalf("expected type %d, got %d", h.Type, buf[0])
	}

	if buf[1] != byte(h.SeqNum) {
		t.Fatalf("expected seqnum %d, got %d", h.SeqNum, buf[1])
	}

	if buf[2] != byte(h.Channel) {
		t.Fatalf("expected channel %d, got %d", h.Channel, buf[2])
	}
}

func TestPayloadBytes(t *testing.T, pb []byte, expected []byte) {
	if len(pb) != len(expected) {
		t.Fatalf("expected %d bytes, got %d", len(expected), len(pb))
	}

	if !reflect.DeepEqual(pb, expected) {
		t.Fatalf("expected %v, got %v", expected, pb)
	}
}

func TestPacketDecode(t *testing.T, b []byte, dp pcw.DecodePacket) pcw.Packet {
	h, pb, err := pcw.ParseHeader(b)
	if err != nil {
		t.Fatal(err)
	}

	p, err := dp(h, pb)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestHeader(t *testing.T, h pcw.Header, exp pcw.Header) {
	if h.Type != exp.Type {
		t.Fatalf("expected type %d, got %d", exp.Type, h.Type)
	}

	if h.SeqNum != exp.SeqNum {
		t.Fatalf("expected seqnum %d, got %d", exp.SeqNum, h.SeqNum)
	}

	if h.Channel != exp.Channel {
		t.Fatalf("expected channel %d, got %d", exp.Channel, h.Channel)
	}
}
