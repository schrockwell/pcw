# Packetized Keying Protocol (PKP) Specification

Version 1

Updated 2023-11-26

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

All numeric values in example blocks are in hexadecimal notation. Brackets are purely notational, to show the groupings of multi-byte values.

All packet types and fields are introduced in specification version 1, unless otherwise noted.

# Definitions

- **Client** - a device that takes Morse input from a user and emits PKP packets
- **Server** - a device that takes in PKP packets and recreates the Morse for RF transmission
- **Timestamp** - a monotonically-increasing uint32 with microsecond precision
- **Timed Packet** - a packet containing a Timestamp field (Key Up, Key Down, or Element packet types)
- **Sync** - a timestamp of 0, indicating that timestamps of upcoming Timed Packets should be measured relative to the time that the Sync was received by the server
- **Channel** - a value from 0 to 127 representing a line to be keyed

# Field Types

All mutli-byte values use big-endian (network) byte order.

- **uint8** - unsigned 8-bit integer (0 to 255)
- **uint16** - unsigned 16-bit integer (0 to 65535)
- **uint32** - unsigned 32-bit integer (0 to 4294967295)

# Packet Structure

The packet has the following structure:

    [header] [payload]

The total length of the packet is not fixed, and depends on both the packet type and payload contents.

An entire PKP packet SHOULD fit within the payload of a single UDP packet. The useful payload of a UDP packet may be determined by taking a holistic look at the entire network path, its minimum MTU, IP and UDP overhead, and additional overhead due to application protocols such as DTLS and SCTP.

# Header

## Header Fields

The header is composed of the following fields:

1. **Header Length (uint8)** - The number of bytes to follow which will comprise the header. Currently this is 5 bytes, but it may increase in future protocol versions as header fields are added.

2. **Payload Length (uint16)** - The number of bytes after the header, which comprise the payload.

3. **Packet Type (uint8)** - One of the predefined packet types, which describes the contents of the payload.

4. **Sequence Number (uint8)** - The sequence number MUST be incremented by one for every packet sent, wrapping around from 255 to 0. The sequence number SHOULD be used to detect if a packet was lost or received out-of-order.

5. **Address (uint8)** - The address allows multiple devices to be controlled by a single server. Address 0 MUST be implemented as the default device. Addresses 1 through 255 MAY be used if additional devices are to be controlled.

## Header Example

Header length 5 bytes, payload length 4 bytes, payload type 1, sequence number 171, address 0

    HEADER                      PAYLOAD
    --------------------------- -------------
    [05] [00 04] [01] [AB] [00] [xx xx xx xx]

# Packet Types

## Key Up (0x00)

Sent by: client only

Releases the remote key.

The server MUST reply with a Dropped packet if the timestamp is too late to be queued up for future playing.

### Key Up Fields

1. **Channel (uint8)**
2. **Timestamp (uint32)**

### Key Up Examples

Sync

    HEADER                      PAYLOAD
    --------------------------- ------------------
    [05] [00 04] [01] [AB] [00] [00] [00 00 00 00]

Timestamp 17965876 µs

    HEADER                      PAYLOAD
    --------------------------- ------------------
    [05] [00 04] [01] [AB] [00] [00] [01 12 23 34]

## Key Down (0x01)

Sent by: client only

Engages the remote key.

The server MUST reply with a Dropped packet if the timestamp is too late to be queued up for future playing.

The server MAY implement a time-out timer (TOT) to forcefully release the key if a Key Up packet is not received after a reasonable duration.

### Key Down Fields

1. **Channel (uint8)**
2. **Timestamp (uint32)**

## Element (0x02)

Sent by: client only

Engages the remote key for a specified duration.

The server MUST reply with a Dropped packet if the timestamp is too late to be queued up for future playing.

### Element Fields

1. **Channel (uint8)**
2. **Timestamp (uint32)**
3. **Duration (uint32)**

## Characters (0x03)

Sent by: client only

Enqueues characters for CW generation at the server.

The characters MUST be encoded as Latin-1 per [ISO 8859-1](https://en.wikipedia.org/wiki/ISO/IEC_8859-1).

The characters MAY be sent to any destination that can recreate the Morse representation of those characters, such as a radio, WinKeyer serial device, or software. Devices MAY interpret characters in special ways, such as for custom spacing or prosigns.

The server MAY ignore the packet if there is no applicable device to recreate the Morse code.

### Characters Fields

1. **Channel (uint8)**
2. **Length (uint8)**
3. **Characters (uint8 array)**

## WinKeyer Command (0x04)

Sent by: client only

Sends commands to WinKeyer device.

The bytes MUST correspond to Host Mode commands as described in the [WinKeyer documentation](https://hamcrafters2.com/files/WK3_Datasheet_v1.3.pdf) (TODO: which version, specifically?).

The server MAY forward these commands directly to a WinKeyer device.

In the absence of a WinKeyer device, the server MAY interpret the commands to control alternative hardware or software keying. For example, the server could handle the WinKey "Set WPM Speed" command to set the WPM of a radio using the radio's command protocol.

The server MAY ignore the packet if there is no applicable device to handle the particular command.

### WinKeyer Command Fields

1. **Channel (uint8)**
2. **Length (uint8)**
3. **Command Data (uint8 array)**

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

Indicates that a Timed Packet arrived too late.

The server MUST send this packet if a Timed Packet was received too late to be recreated.

## Application Data (0x09)

Sent by: client or server

Sends application-specific data.

The client or server MAY use this packet type to send any custom payload for control or informational purposes that are beyond the scope of the PKP protocol.

## Undefined (0xTODO through 0xFF)

The device MUST ignore any packet type that it does not implement.

# Backwards Compatibility

Protocol changes SHALL NOT modify the meaning of any existing packet types or fields, and SHALL NOT change the ordering or encoding of existing fields.

Protocol changes MAY add new packet types, and MAY add new fields to existing packet types. New fields SHALL always be appended such that their bytes occur after any existing fields. Therefore, the lengths of headers and payloads MAY increase in new protocol versions.

In the event that the header length or payload length is longer than the device expects, this indicates the packet contains information from a newer version of the protocol. The device MUST accept the packet and MUST ignore the additional bytes.

# Serializing Packets

TODO

    PREAMBLE      PACKET             TRAILER
    -----------   ------------------ ----------
    [AA AA AA AA] [header] [payload] [checksum]

- Preamble
  - Magic delimiter (uint32) - the value `0xAAAAAAAA`
- Packet - the original packet header and payload
- Trailer
  - Checksum (uint8) - the sum of all the packet bytes (modulo 256)
