package share

const (
	// ExamplePacketSerialkey is Serialkey(uint32) + PacketSize(uint16)
	ExamplePacketSerialkey uint32 = 0x5da9c31b
	ExamVer                uint32 = 10
)

const (
	// ResultSuccess is Packet Request Success
	ResultSuccess uint32 = 0x00
	// ResultFail is Packet Request Fail
	ResultFail uint32 = 0x01
)

const (
	packetCommandStart uint32 = 0x100
	// C2SPacketCommandLoginUserReq is
	C2SPacketCommandLoginUserReq = packetCommandStart + 1
	// S2CPacketCommandLoginUserRes is
	S2CPacketCommandLoginUserRes = packetCommandStart + 2

	// C2SPacketCommandGolobalMsgReq is
	C2SPacketCommandGolobalMsgReq = packetCommandStart + 3
	// S2CPacketCommandLobbyMsgRes is
	S2CPacketCommandLobbyMsgRes = packetCommandStart + 4

	// C2SPacketCommandRoomEnterReq is
	C2SPacketCommandRoomEnterReq = packetCommandStart + 5
	// S2CPacketCommandRoomEnterRes is
	S2CPacketCommandRoomEnterRes = packetCommandStart + 6

	// C2SPacketCommandRoomLeaveReq is
	C2SPacketCommandRoomLeaveReq = packetCommandStart + 7
	// S2CPacketCommandRoomLeaveRes is
	S2CPacketCommandRoomLeaveRes = packetCommandStart + 8

	// C2SPacketCommandRoomMsgReq is
	C2SPacketCommandRoomMsgReq = packetCommandStart + 9
	// S2CPacketCommandRoomMsgRes is
	S2CPacketCommandRoomMsgRes = packetCommandStart + 10
)

// C2SPCLoginUserReq is  C2SPacketCommandLoginUserReq packet struct
type C2SPCLoginUserReq struct {
	UserID string
}

// S2CPCLoginUserRes is  S2CPacketCommandLoginUserRes packet struct
type S2CPCLoginUserRes struct {
	Result uint32
	UserSn uint32
	UserID string
}

// C2SPCLobbySendMsgReq is  C2SPacketCommandLoginUserReq packet struct
type C2SPCLobbySendMsgReq struct {
	UserSn uint32
	Msg    string
}

// S2CPCLobbySendMsgRes is  S2CPacketCommandLoginUserRes packet struct
type S2CPCLobbySendMsgRes struct {
	Result uint32
	Userid string
	Msg    string
}

// C2SPCRoomEnterReq is  C2SPacketCommandRoomEnterReq packet struct
type C2SPCRoomEnterReq struct {
	UserSn   uint32
	RoomName string
}

// S2CPCRoomEnterRes is  C2SPacketCommandRoomLeaveReq packet struct
type S2CPCRoomEnterRes struct {
	Result  uint32
	roomIdx uint32
}

// C2SPCRoomLeaveReq is  C2SPacketCommandRoomEnterReq packet struct
type C2SPCRoomLeaveReq struct {
	UserSn   uint32
	RoomName string
}

// S2CPCRoomLeaveRes is  S2CPacketCommandRoomLeaveRes packet struct
type S2CPCRoomLeaveRes struct {
	Result uint32
	Userid string
	Msg    string
}
