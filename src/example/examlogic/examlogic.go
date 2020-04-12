package examlogic

const (
	// ExamplePacketSerialkey is Serialkey(uint32) + PacketSize(uint16)
	ExamplePacketSerialkey uint32 = 0x5da9c31b
)

const (
	packetCommandStart uint32 = 0x100
	// C2SPacketCommandLoginUserReq is
	C2SPacketCommandLoginUserReq = packetCommandStart + 1
	// S2CPacketCommandLoginUserRes is
	S2CPacketCommandLoginUserRes = packetCommandStart + 2

	// C2SPacketCommandLoginUserReq is
	C2SPacketCommandGolobalSendMsgReq = packetCommandStart + 3
	// S2CPacketCommandLoginUserRes is
	S2CPacketCommandGolobalSendMsgRes = packetCommandStart + 4
)
