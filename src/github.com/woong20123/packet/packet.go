package packet

import (
	"encoding/binary"
)

const (
	// PacketHeaderSize is Serialkey(uint32) + PacketSize(uint16)
	PacketHeaderSize uint32 = 10
	// MaxPacketSize is Serialkey(uint32) + PacketSize(uint16)
	MaxPacketSize uint16 = 0x1000
)

// HeaderChack check if the packet header normal
func HeaderChack(buffer []byte, serialkey uint32) bool {
	if serialkey == binary.LittleEndian.Uint32(buffer) {
		PacketSize := binary.LittleEndian.Uint16(buffer[4:])
		if PacketSize <= MaxPacketSize {
			return true
		}
	}
	return false
}

// Header is
type Header struct {
	serialkey     uint32
	packetSize    uint16
	packetCommand uint32
}

// Packet is
type Packet struct {
	header Header
	buffer []byte
}

// NewPacket allocate memory in the packet buffer
func NewPacket(buffersize uint32) *Packet {
	p := Packet{}
	p.header = Header{}
	p.buffer = make([]byte, buffersize)
	return &p
}

// GetByte packet transform to []byte
func (p *Packet) GetByte() []byte {
	buffer := make([]byte, PacketHeaderSize+uint32(p.getSize()))
	binary.LittleEndian.PutUint32(buffer, p.header.serialkey)
	binary.LittleEndian.PutUint16(buffer[4:], p.header.packetSize)
	binary.LittleEndian.PutUint32(buffer[6:], p.header.packetCommand)
	copy(buffer[10:], p.buffer[:p.header.packetSize])
	return buffer
}

// GetCommand is return packet command
func (p *Packet) GetCommand() uint32 {
	return p.header.packetCommand
}

// SetHeader Set Header Info
func (p *Packet) SetHeader(serialKey uint32, packetsize uint16, packetcommand uint32) {
	p.header.serialkey = serialKey
	p.header.packetSize = packetsize
	p.header.packetCommand = packetcommand
}

// CopyByte write byte to packet buffer
func (p *Packet) CopyByte(data []byte) {
	length := copy(p.buffer[p.getSize():], data)
	p.addSize(uint16(length))
}

// Write write various type to packet buffer
func (p *Packet) Write(data ...interface{}) {
	for _, value := range data {
		switch v := value.(type) {
		case uint8:
		case int8:
			p.buffer[p.getSize()] = byte(v)
			p.addSize(1)
		case uint16:
		case int16:
			binary.LittleEndian.PutUint16(p.buffer[p.getSize():], uint16(v))
			p.addSize(2)
		case uint32:
		case int32:
		case int:
		case float32:
			binary.LittleEndian.PutUint32(p.buffer[p.getSize():], uint32(v))
			p.addSize(4)
		case uint64:
		case int64:
		case float64:
			binary.LittleEndian.PutUint64(p.buffer[p.getSize():], uint64(v))
			p.addSize(8)
		case string:
			b := []byte(v)
			b = append(b, 0)
			length := copy(p.buffer[p.getSize():], b)
			p.addSize(uint16(length))
		}
	}
}

func (p *Packet) getSize() uint16 {
	return p.header.packetSize
}

func (p *Packet) setSize(size uint16) {
	p.header.packetSize = size
}

func (p *Packet) addSize(size uint16) {
	p.header.packetSize += size
}
