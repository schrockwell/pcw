# Packetized CW (PCW) Protocol

This is a work in progress! Nothing is finalized yet.

## Goals

1. **Three levels of service**

   1. Per-transition: In high-reliability, low-jitter network situations, like on a local LAN, PCW can carry individual key-up and key-down transitions for the most accurate and lowest-latency recreation of a CW signal.

   2. Per-element: In a less-optimal network, the user may send entire elements (dits and dahs) instead.

   3. Per-character: If precise timing is not needed, the user may send individual characters for recreation at the server end.

2. **Simple encoding** - PCW is designed to be easy to encode and decode, requiring only basic byte manipulation, which is readily achievable by any microcontroller or software.

3. **Small** - Packets are extremely compact and should fit within the payload of a single UDP packet.

## Discussion

### Latency and jitter

### Packet reordering

### Packet loss

## Specification

See [the specification](SPECIFICATION.md).
