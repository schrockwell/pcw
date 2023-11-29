# Packetized Keying Protocol (PKP)

This is a work in progress! Nothing is finalized yet.

## Goals

1. **Two levels of service**

   1. Per-transition: PKP can carry individual key-up and key-down transitions for the most accurate recreation of a keying signal.

   2. Per-character: If precise timing is not needed, the user may send individual characters for recreation at the server end.

2. **Simple encoding** - PKP is designed to be easy to encode and decode, requiring only basic byte manipulation, which is readily achievable by any software or microcontroller.

3. **Small** - Packets are extremely compact, designed to fit within the payload of a single UDP packet.

4. **Not just for CW** - Keying applies to many aspects of radio beyond Morse code. For example, push-to-talk (PTT), radioteletype (RTTY), and generic relay controls are also supported use-cases.

5. **Backwards-compatible** - Additions can be made to the protocol without breaking existing implementations.

## Discussion

### Latency and jitter

### Packet reordering

### Packet loss

## Specification

See [the specification](SPECIFICATION.md).
