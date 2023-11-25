package pcw

import "errors"

var ErrInvalidHeader = errors.New("invalid header")
var ErrInvalidPacket = errors.New("invalid packet")
var ErrInvalidPacketType = errors.New("invalid packet type")
