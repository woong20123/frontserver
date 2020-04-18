package packet

import (
	"encoding/binary"
	"log"
	"math"
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
	header  Header
	buffer  []byte
	readPos uint16
}

// NewPacket allocate memory in the packet buffer
func NewPacket() *Packet {
	p := Packet{}
	p.header = Header{}
	p.buffer = make([]byte, MaxPacketSize)
	p.readPos = 0
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
func (p *Packet) Write(datas ...interface{}) {
	for _, data := range datas {
		if t, _ := intDataSize(data); t != 0 {
			typeSize := uint16(t)
			switch v := data.(type) {
			case *bool:
				if *v {
					p.buffer[p.getSize()] = 1
				} else {
					p.buffer[p.getSize()] = 0
				}
				p.addSize(typeSize)
			case bool:
				if v {
					p.buffer[p.getSize()] = 1
				} else {
					p.buffer[p.getSize()] = 0
				}
				p.addSize(typeSize)
			case []bool:
				for _, x := range v {
					if x {
						p.buffer[p.getSize()] = 1
					} else {
						p.buffer[p.getSize()] = 0
					}
					p.addSize(typeSize)
				}
			case *int8:
				p.buffer[p.getSize()] = byte(*v)
				p.addSize(typeSize)
			case int8:
				p.buffer[p.getSize()] = byte(v)
				p.addSize(typeSize)
			case []int8:
				for _, x := range v {
					p.buffer[p.getSize()] = byte(x)
					p.addSize(typeSize)
				}
			case *uint8:
				p.buffer[p.getSize()] = *v
				p.addSize(typeSize)
			case uint8:
				p.buffer[p.getSize()] = v
				p.addSize(typeSize)
			case []uint8:
				for _, x := range v {
					p.buffer[p.getSize()] = x
					p.addSize(typeSize)
				}
			case *int16:
				binary.LittleEndian.PutUint16(p.buffer[p.getSize():], uint16(*v))
				p.addSize(typeSize)
			case int16:
				binary.LittleEndian.PutUint16(p.buffer[p.getSize():], uint16(v))
				p.addSize(typeSize)
			case []int16:
				for _, x := range v {
					binary.LittleEndian.PutUint16(p.buffer[p.getSize():], uint16(x))
					p.addSize(typeSize)
				}
			case *uint16:
				binary.LittleEndian.PutUint16(p.buffer[p.getSize():], *v)
				p.addSize(typeSize)
			case uint16:
				binary.LittleEndian.PutUint16(p.buffer[p.getSize():], v)
				p.addSize(typeSize)
			case []uint16:
				for _, x := range v {
					binary.LittleEndian.PutUint16(p.buffer[p.getSize():], x)
					p.addSize(typeSize)
				}
			case *int:
			case *int32:
				binary.LittleEndian.PutUint32(p.buffer[p.getSize():], uint32(*v))
				p.addSize(typeSize)
			case int:
			case int32:
				binary.LittleEndian.PutUint32(p.buffer[p.getSize():], uint32(v))
				p.addSize(typeSize)
			case []int32:
				for _, x := range v {
					binary.LittleEndian.PutUint32(p.buffer[p.getSize():], uint32(x))
					p.addSize(typeSize)
				}
			case *uint:
			case *uint32:
				binary.LittleEndian.PutUint32(p.buffer[p.getSize():], *v)
				p.addSize(typeSize)
			case uint:
			case uint32:
				binary.LittleEndian.PutUint32(p.buffer[p.getSize():], v)
				p.addSize(typeSize)
			case []uint32:
				for _, x := range v {
					binary.LittleEndian.PutUint32(p.buffer[p.getSize():], x)
					p.addSize(typeSize)
				}
			case *int64:
				binary.LittleEndian.PutUint64(p.buffer[p.getSize():], uint64(*v))
				p.addSize(typeSize)
			case int64:
				binary.LittleEndian.PutUint64(p.buffer[p.getSize():], uint64(v))
				p.addSize(typeSize)
			case []int64:
				for _, x := range v {
					binary.LittleEndian.PutUint64(p.buffer[p.getSize():], uint64(x))
					p.addSize(typeSize)
				}
			case *uint64:
				binary.LittleEndian.PutUint64(p.buffer[p.getSize():], *v)
				p.addSize(typeSize)
			case uint64:
				binary.LittleEndian.PutUint64(p.buffer[p.getSize():], v)
				p.addSize(typeSize)
			case []uint64:
				for _, x := range v {
					binary.LittleEndian.PutUint64(p.buffer[p.getSize():], x)
				}
			case *float32:
				binary.LittleEndian.PutUint32(p.buffer[p.getSize():], math.Float32bits(*v))
				p.addSize(typeSize)
			case float32:
				binary.LittleEndian.PutUint32(p.buffer[p.getSize():], math.Float32bits(v))
				p.addSize(typeSize)
			case []float32:
				for _, x := range v {
					binary.LittleEndian.PutUint32(p.buffer[p.getSize():], math.Float32bits(x))
					p.addSize(typeSize)
				}
			case *float64:
				binary.LittleEndian.PutUint64(p.buffer[p.getSize():], math.Float64bits(*v))
				p.addSize(typeSize)
			case float64:
				binary.LittleEndian.PutUint64(p.buffer[p.getSize():], math.Float64bits(v))
				p.addSize(typeSize)
			case []float64:
				for _, x := range v {
					binary.LittleEndian.PutUint64(p.buffer[p.getSize():], math.Float64bits(x))
					p.addSize(typeSize)
				}
			case string:
				length := uint16(len(v))
				binary.LittleEndian.PutUint16(p.buffer[p.getSize():], length)
				p.addSize(2)
				len := copy(p.buffer[p.getSize():], v)
				p.addSize(uint16(len))
			case *string:
				length := uint16(len(*v))
				binary.LittleEndian.PutUint16(p.buffer[p.getSize():], length)
				p.addSize(2)
				len := copy(p.buffer[p.getSize():], *v)
				p.addSize(uint16(len))
			default:
				log.Printf("not find type Packet Write")
			}
		} else {
			log.Printf("error Packet Write")
		}
	}
}

// WriteString is
// func (p *Packet) WriteString(str string) {
// 	length := uint16(len(str))
// 	p.Write(length)
// 	len := copy(p.buffer[p.getSize():], str)
// 	p.addSize(uint16(len))
// }

func (p *Packet) Read(datas ...interface{}) {
	for _, data := range datas {
		if t, _ := intDataSize(data); t != 0 {
			typeSize := uint16(t)
			switch data := data.(type) {
			case *bool:
				*data = p.buffer[p.getReadPos()] != 0
				p.addReadPos(typeSize)
			case *int8:
				*data = int8(p.buffer[p.getReadPos()])
				p.addReadPos(typeSize)
			case *uint8:
				*data = p.buffer[p.getReadPos()]
				p.addReadPos(typeSize)
			case *int16:
				*data = int16(binary.LittleEndian.Uint16(p.buffer[p.getReadPos():]))
				p.addReadPos(typeSize)
			case *uint16:
				*data = binary.LittleEndian.Uint16(p.buffer[p.getReadPos():])
				p.addReadPos(typeSize)
			case *int32:
				*data = int32(binary.LittleEndian.Uint32(p.buffer[p.getReadPos():]))
				p.addReadPos(typeSize)
			case *int:
				*data = int(binary.LittleEndian.Uint32(p.buffer[p.getReadPos():]))
				p.addReadPos(typeSize)
			case *uint32:
				*data = binary.LittleEndian.Uint32(p.buffer[p.getReadPos():])
				p.addReadPos(typeSize)
			case *uint:
				*data = uint(binary.LittleEndian.Uint32(p.buffer[p.getReadPos():]))
				p.addReadPos(typeSize)
			case *int64:
				*data = int64(binary.LittleEndian.Uint64(p.buffer[p.getReadPos():]))
				p.addReadPos(typeSize)
			case *uint64:
				*data = binary.LittleEndian.Uint64(p.buffer[p.getReadPos():])
				p.addReadPos(typeSize)
			case *float32:
				*data = math.Float32frombits(binary.LittleEndian.Uint32(p.buffer[p.getReadPos():]))
				p.addReadPos(typeSize)
			case *float64:
				*data = math.Float64frombits(binary.LittleEndian.Uint64(p.buffer[p.getReadPos():]))
				p.addReadPos(typeSize)
			case []bool:
				for i := range data {
					data[i] = p.buffer[p.getReadPos()] != 0
					p.addReadPos(typeSize)
				}
			case []int8:
				for i := range data {
					data[i] = int8(p.buffer[p.getReadPos()])
					p.addReadPos(typeSize)
				}
			case []uint8:
				for i := range data {
					data[i] = uint8(p.buffer[p.getReadPos()])
					p.addReadPos(typeSize)
				}
			case []int16:
				for i := range data {
					data[i] = int16(binary.LittleEndian.Uint16(p.buffer[p.getReadPos():]))
					p.addReadPos(typeSize)
				}
			case []uint16:
				for i := range data {
					data[i] = binary.LittleEndian.Uint16(p.buffer[p.getReadPos():])
					p.addReadPos(typeSize)
				}
			case []int32:
				for i := range data {
					data[i] = int32(binary.LittleEndian.Uint32(p.buffer[p.getReadPos():]))
					p.addReadPos(typeSize)
				}
			case []uint32:
				for i := range data {
					data[i] = binary.LittleEndian.Uint32(p.buffer[p.getReadPos():])
					p.addReadPos(typeSize)
				}
			case []int64:
				for i := range data {
					data[i] = int64(binary.LittleEndian.Uint64(p.buffer[p.getReadPos():]))
					p.addReadPos(typeSize)
				}
			case []uint64:
				for i := range data {
					data[i] = binary.LittleEndian.Uint64(p.buffer[p.getReadPos():])
					p.addReadPos(typeSize)
				}
			case []float32:
				for i := range data {
					data[i] = math.Float32frombits(binary.LittleEndian.Uint32(p.buffer[p.getReadPos():]))
					p.addReadPos(typeSize)
				}
			case []float64:
				for i := range data {
					data[i] = math.Float64frombits(binary.LittleEndian.Uint64(p.buffer[p.getReadPos():]))
				}
			case *string:
				length := binary.LittleEndian.Uint16(p.buffer[p.getReadPos():])
				p.addReadPos(2)
				bs := make([]byte, length)
				len := copy(bs, p.buffer[p.getReadPos():p.getReadPos()+length])
				p.addReadPos(uint16(len))
				*data = string(bs)
			default:
				log.Printf("not find tyep Packet Read")
			}
		} else {
			log.Printf("error Packet Read")
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

func (p *Packet) getReadPos() uint16 {
	return p.readPos
}

func (p *Packet) setReadPos(pos uint16) {
	p.readPos = pos
}

func (p *Packet) addReadPos(pos uint16) {
	p.readPos += pos
}

// AssemblyFromBuffer is make packet from buffer
func AssemblyFromBuffer(buffer []byte, bufferpos uint32, serialkey uint32) (resultpos uint32, pPacket *Packet) {
	resultpos = bufferpos
	var headerFind bool = false
	if resultpos < PacketHeaderSize {
		return
	}

	findIndex := resultpos - PacketHeaderSize

	// 전달받은 버퍼를 순회하면서 패킷이 정상적으로 전달되었는지 확인합니다.
	for index := range buffer {
		// must receive At least as much as the header
		// 최소한 헤더만큼 수신받아야 처리 가능,
		// 이곳에 들어왔다면 검사한 버퍼를 모두 버립니다.
		if uint32(index) > findIndex {
			copy(buffer, buffer[index:resultpos])
			resultpos = resultpos - uint32(index)
			break
		}

		if true == HeaderChack(buffer[index:], serialkey) {
			headerFind = true

			// 만약 패킷 헤더가 처음이 아니라면 나머지 버퍼를 버립니다.
			if index != 0 {
				copy(buffer, buffer[index:resultpos])
				resultpos = resultpos - uint32(index)
			}
			break
		}
	}

	// 패킷 시작지점을 찾았다면
	if true == headerFind {
		serialKey := uint32(binary.LittleEndian.Uint32(buffer))
		packetSize := uint32(binary.LittleEndian.Uint16(buffer[4:]))
		packetCommand := uint32(binary.LittleEndian.Uint16(buffer[6:]))
		totalPacketSize := packetSize + PacketHeaderSize

		// 패킷을 만들 수 있을 만큼 패킷을 전달 받았다면 패킷을 만들고
		// Logic 처리 goroutine에 전달합니다.
		if totalPacketSize <= resultpos {
			pPacket = NewPacket()
			pPacket.SetHeader(serialKey, 0, packetCommand)
			pPacket.CopyByte(buffer[PacketHeaderSize:totalPacketSize])

			copy(buffer, buffer[totalPacketSize:resultpos])
			resultpos = resultpos - totalPacketSize
		}
	}
	return
}

// intDataSize returns the size of the data required to represent the data when encoded.
// It returns zero if the type cannot be implemented by the fast path in Read or Write.
func intDataSize(data interface{}) (typeSize int, TotalSize int) {
	typeSize = 0
	TotalSize = 0
	switch data := data.(type) {
	case bool, int8, uint8, *bool, *int8, *uint8:
		typeSize = 1
		TotalSize = 1
	case []bool:
		typeSize = 1
		TotalSize = len(data)
	case []int8:
		typeSize = 1
		TotalSize = len(data)
	case []uint8:
		typeSize = 1
		TotalSize = len(data)
	case int16, uint16, *int16, *uint16:
		typeSize = 2
		TotalSize = 2
	case []int16:
		typeSize = 2
		TotalSize = 2 * len(data)
	case []uint16:
		typeSize = 2
		TotalSize = 2 * len(data)
	case int32, uint32, *int32, *uint32:
		typeSize = 4
		TotalSize = 4
	case []int32:
		typeSize = 4
		TotalSize = 4 * len(data)
	case []uint32:
		typeSize = 4
		TotalSize = 4 * len(data)
	case int64, uint64, *int64, *uint64:
		typeSize = 8
		TotalSize = 8
	case []int64:
		typeSize = 8
		TotalSize = 8 * len(data)
	case []uint64:
		typeSize = 8
		TotalSize = 8 * len(data)
	case float32, *float32:
		typeSize = 4
		TotalSize = 4
	case float64, *float64:
		typeSize = 8
		TotalSize = 8
	case []float32:
		typeSize = 4
		TotalSize = 4 * len(data)
	case []float64:
		typeSize = 8
		TotalSize = 8 * len(data)
	case string, *string:
		typeSize = 4
		TotalSize = 4
	}
	return
}
