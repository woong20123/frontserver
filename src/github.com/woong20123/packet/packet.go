package packet

import (
	"bytes"
	"encoding/binary"
	"log"
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
	buf     *bytes.Buffer
	readPos uint16
}

// NewPacket allocate memory in the packet buffer
func NewPacket() *Packet {
	p := Packet{}
	p.header = Header{}
	p.buf = new(bytes.Buffer)
	p.readPos = 0
	return &p
}

// GetByteBuf is
func (p *Packet) GetByteBuf() *bytes.Buffer {
	return p.buf
}

// GetByte packet transform to []byte
func (p *Packet) GetByte() []byte {
	buffer := make([]byte, PacketHeaderSize+uint32(p.getSize()))
	binary.LittleEndian.PutUint32(buffer, p.header.serialkey)
	binary.LittleEndian.PutUint16(buffer[4:], p.header.packetSize)
	binary.LittleEndian.PutUint32(buffer[6:], p.header.packetCommand)
	copy(buffer[10:], p.buf.Bytes())
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
	p.buf.Write(data)
	p.addSize(uint16(len(data)))
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

func (p *Packet) addReadPos(size uint16) {
	p.readPos += size
}

func (p *Packet) Write(order binary.ByteOrder, data interface{}) {
	err := binary.Write(p.buf, order, data)
	if err != nil {
		log.Println("Packet Write Fail", err)
		return
	}
	p.setSize(uint16(p.buf.Len()))
}

// WriteString is
func (p *Packet) WriteString(order binary.ByteOrder, data string) {
	var length uint16 = uint16(len(data))
	err := binary.Write(p.buf, order, length)
	if err != nil {
		log.Println("StringWrite Fail", err)
		return
	}
	_, err = p.buf.WriteString(data)
	if err != nil {
		log.Println("StringWrite Fail", err)
		return
	}
	p.setSize(uint16(p.buf.Len()))
}

func (p *Packet) Read(order binary.ByteOrder, data interface{}) {
	err := binary.Read(p.buf, binary.LittleEndian, &data)
	if err != nil {
		log.Println("Packet Read Fail", err)
		return
	}
}

func (p *Packet) ReadString(order binary.ByteOrder, data *string) {
	//err := binary.Read(p.buf, binary.LittleEndian, &data)
	// *data, err = p.buf.ReadString()
	// if err != nil {
	// 	log.Println("Packet Read Fail", err)
	// 	return
	// }
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
			pPacket.SetHeader(serialKey, uint16(packetSize), packetCommand)
			pPacket.CopyByte(buffer[PacketHeaderSize:totalPacketSize])

			copy(buffer, buffer[totalPacketSize:resultpos])
			resultpos = resultpos - totalPacketSize
		}
	}
	return
}
