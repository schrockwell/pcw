# Packetized CW (PCW) Specification

Version: 0.1.0

Updated: 2023-11-26

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

All numeric values in example blocks are in hexadecimal notation. Brackets are purely notational, to show the groupings of multi-byte values.

# Definitions

- **Client** - a device that takes Morse input from a user and emits PCW packets
- **Server** - a device that takes in PCW packets and recreates the Morse for RF transmission
- **Timestamp** - a monotonically-increasing uint32 with microsecond precision
- **Timely Packet** - a packet containing a Timestamp field (Key Up, Key Down, or Element packet types)
- **Sync** - a timestamp of 0, indicating that timestamps of upcoming Timely Packets should be measured relative to the time that the Sync was received by the server

# Field Types

All mutli-byte values use big-endian (network) byte order.

- **uint8** - unsigned 8-bit integer (0 to 255)
- **uint16** - unsigned 16-bit integer (0 to 65535)
- **uint32** - unsigned 32-bit integer (0 to 4294967295)

# Packet Structure

The packet has the following structure:

    [header] [payload]

The total length of the packet is not fixed, and depends on both the packet type and payload contents.

An entire PCW packet SHOULD fit within the payload of a single UDP packet. The useful payload of a UDP packet may be determined by taking a holistic look at the entire network path, its minimum MTU, IP and UDP overhead, and additional overhead due to application protocols such as DTLS and SCTP.

# Header

## Header Fields

The header is composed of the following fields:

1. **Packet Type (uint8)** - The packet type MUST be one of the defined packet types below.

2. **Sequence Number (uint8)** - The sequence number MUST be incremented by one for every packet sent, wrapping around from 255 to 0. The sequence number SHOULD be used to detect if a packet was lost or received out-of-order.

3. **Channel (uint8)** - The channel allows multiple keys to be controlled by a single server. Channel 0 MUST be implemented as the default channel. Channels 1 through 255 MAY be used if additional keys are supported.

## Header Example

Payload type 1, sequence number 171, channel 0

    [01] [AB] [00]

# Packet Types

## Key Up (0x00)

Sent by: client only

Releases the remote key.

The server MUST reply with a Dropped packet if the timestamp is too late to be queued up for future playing.

### Key Up Fields

1. **Timestamp (uint32)** - A time, in microseconds, when to release the key. The timestamp of 0 MUST be reserved as a synchronization signal. The server and client MUST support accurate timing when the timestamp wraps around from UINT32_MAX to 0.

### Key Up Examples

Payload type 0, sequence number 4, channel 0, sync

    [00] [04] [00] [00 00 00 00]

Payload type 0, sequence number 10, channel 0, timestamp 17965876 µs

    [00] [12] [00] [01 12 23 34]

## Key Down (0x01)

Sent by: client only

Engages the remote key.

The server MUST reply with a Dropped packet if the timestamp is too late to be queued up for future playing.

The server MAY implement a time-out timer (TOT) to forcefully release the key if a Key Up packet is not received after a reasonable duration.

### Key Down Fields

1. Timestamp (uint32)  
   A time, in microseconds, when to trigger the key. The timestamp of 0 MUST be reserved as a synchronization signal. The server and client MUST support accurate timing when the timestamp wraps around from UINT32_MAX to 0.

### Key Down Examples

Payload type 1, sequence number 4, channel 0, sync

    [01] [04] [00] [00 00 00 00]

Payload type 1, sequence number 10, channel 0, timestamp 17965876 µs

    [01] [12] [00] [01 12 23 34]

## Element (0x02)

Sent by: client only

Engages the remote key for a specified duration.

The server MUST reply with a Dropped packet if the timestamp is too late to be queued up for future playing.

## Characters (0x03)

Sent by: client only

Enqueues characters for CW generation at the server.

The characters MUST be encoded as ASCII (TODO: which ISO??).

The characters MAY be sent to any destination that can recreate the Morse representation of those characters, such as a radio, WinKeyer-compatible serial device, or software. Devices MAY interpret characters in special ways, such as for custom spacing or prosigns.

The server MAY ignore the packet if there is no applicable device to recreate the Morse code.

## WinKeyer Command (0x04)

Sent by: client only

Sends commands to WinKeyer-compatible device.

The bytes MUST correspond to Host Mode commands as described in the [WinKeyer documentation](https://hamcrafters2.com/files/WK3_Datasheet_v1.3.pdf) (TODO: which version, specifically?).

The server MAY forward these commands directly to a WinKeyer-compatible device.

In the absence of a WinKeyer-compatible device, the server MAY introspect the commands to control alternative hardware or software keying. For example, the server could handle the WinKey "Set WPM Speed" command to set the WPM of a remote radio using its native command protocol.

The server MAY ignore the packet if there is no applicable device to handle the particular command.

## WinKeyer Status (TODO)

Sent by: server only

TODO

## Ping (0x05)

Send by: client only

Starts a RTT measurement.

The client MAY occasionally send Ping packets.

The client SHOULD NOT send a Ping if it has not received a Pong within a short interval (interval determined by implementation).

The client MAY send a Ping if it has not received a Pong after a long interval (interval determined by implementation).

The server MUST respond with a Pong packet having the same timestamp as the Ping packet.

## Pong (0x05)

Send by: server only

Responds to a RTT measurement.

## Missed (0x07)

Sent by: server only

Indicates that a packet was never received.

The server MUST send this packet if it detects a missing client packet based on the sequence number. The threshold of time after which a packet is considered "missed" is left up to the server implementation.

The client SHOULD re-send the missed packet.

## Dropped (0x08)

Sent by: server only

Indicates that a Timely Packet arrived too late.

The server MUST send this packet if a Timely Packet was received too late to be recreated.

## Application Data (0x09)

Sent by: client or server

Sends application-specific data.

The client or server MAY use this packet type to send any custom payload for control or informational purposes that are beyond the scope of the PCW protocol.

# Serializing Packets

TODO

    PREAMBLE               PACKET             TRAILER
    ---------------------- ------------------ ----------
    [AA AA AA AA] [length] [header] [payload] [checksum]

- Preamble
  - Magic delimiter (uint32) - the value `0xAAAAAAAA`
  - Length (uint16) - the number of bytes in the packet to follow, not including the trailer
- Packet - the original packet header and payload
- Trailer
  - Checksum (uint8) - the sum of all the packet bytes (modulo 256)
