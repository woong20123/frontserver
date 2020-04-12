package packet

import (
	"encoding/binary"
)

const (
	// PacketHeaderSize is Serialkey(uint32) + PacketSize(uint16)
	PacketHeaderSize uint32 = 6
	// ExamplePacketSerialkey is Serialkey(uint32) + PacketSize(uint16)
	ExamplePacketSerialkey uint32 = 0x5da9c31b
	// MaxPacketSize is Serialkey(uint32) + PacketSize(uint16)
	MaxPacketSize uint16 = 0x1000
)

// HeaderChack check if the packet header normal
func HeaderChack(buffer []byte) bool {
	if ExamplePacketSerialkey == binary.LittleEndian.Uint32(buffer) {
		PacketSize := binary.LittleEndian.Uint16(buffer[4:])
		if PacketSize <= MaxPacketSize {
			return true
		}
	}
	return false
}

// Header is
type Header struct {
	serialkey  uint32
	packetSize uint16
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

// SetHeader Set Header Info
func (p *Packet) SetHeader(serialKey uint32, packetsize uint16) {
	p.header.serialkey = serialKey
	p.header.packetSize = packetsize
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
			p.addSize(uint16(length) + 1)
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
