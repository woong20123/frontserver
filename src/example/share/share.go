package share

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

	// C2SPacketCommandGolobalSendMsgReq is
	C2SPacketCommandGolobalSendMsgReq = packetCommandStart + 3
	// S2CPacketCommandGolobalSendMsgRes is
	S2CPacketCommandGolobalSendMsgRes = packetCommandStart + 4
)

// C2SPCLoginUserReq is  C2SPacketCommandLoginUserReq packet struct
type C2SPCLoginUserReq struct {
	Userid string
}

// S2CPCLoginUserRes is  S2CPacketCommandLoginUserRes packet struct
type S2CPCLoginUserRes struct {
	Result uint32
}

// C2SPCGolobalSendMsgReq is  C2SPacketCommandLoginUserReq packet struct
type C2SPCGolobalSendMsgReq struct {
	Msg string
}

// S2CPCGolobalSendMsgRes is  S2CPacketCommandLoginUserRes packet struct
type S2CPCGolobalSendMsgRes struct {
	Result uint32
}
