package packet

// Packet is
type PacketPool struct {
	packetpool chan *Packet
}

var instance *PacketPool = nil

// GetInstance is return SingletonObj
func GetPool() *PacketPool {
	if instance == nil {
		instance = new(PacketPool)
		instance.Intialize(2048)
	}
	return instance
}

// newPacket allocate memory in the packet buffer
func newPacket(bAllocPool bool) *Packet {
	p := Packet{}
	p.header = Header{}
	p.buffer = make([]byte, MaxPacketSize)
	p.Init()
	p.bAllocPool = bAllocPool
	return &p
}

// Intialize is
func (pp *PacketPool) Intialize(poolCnt int) {
	pp.packetpool = make(chan *Packet, poolCnt)
	for i := 0; i < poolCnt; i++ {
		pp.packetpool <- newPacket(true)
	}
}

// AcquirePacket is
func (pp *PacketPool) AcquirePacket() *Packet {
	var p *Packet = nil
	select {
	case p = <-pp.packetpool:
	default:
		p = newPacket(false)
	}
	p.Init()
	return p
}

// ReleasePacket is
func (pp *PacketPool) ReleasePacket(p *Packet) {
	if p.bAllocPool {
		pp.packetpool <- p
	}
}
