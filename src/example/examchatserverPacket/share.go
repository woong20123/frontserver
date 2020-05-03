package examchatserverPacket

const (
	// ExamplePacketSerialkey is Serialkey(uint32) + PacketSize(uint16)
	ExamplePacketSerialkey uint32 = 0x5da9c31b
	// ExamVer is exam project version
	ExamVer uint32 = 18
)

const (
	// TCPCliToSvrIdxChat is
	TCPCliToSvrIdxChat uint32 = iota + 0x10
)

const (
	// ResultSuccess is Packet Request Success
	ResultSuccess uint32 = iota
	// ResultFail is Packet Request Fail
	ResultFail
	// ResultExistUserID is Exist UserID
	ResultExistUserID
	// ResultUserStateErr is user's state error
	ResultUserStateErr
	// ResultRoomCreateFail is fail create room object
	ResultRoomCreateFail
	// ResultRoomAlreadyExist is room already exists.
	ResultRoomAlreadyExist
)

const (
	packetSystemCommandStart uint32 = iota + 0x00
	// F2SPacketCmdSysServerRegistReq is
	F2SPacketCmdSysServerRegistReq
	// S2FPacketCmdSysServerRegistRes is
	S2FPacketCmdSysServerRegistRes
)

const (
	packetLogicCommandStart uint32 = iota + 0x2000
	// C2SPacketCommandLoginUserReq is
	C2SPacketCommandLoginUserReq
	// S2CPacketCommandLoginUserRes is
	S2CPacketCommandLoginUserRes

	// C2SPacketCommandLogOutUserReq is
	C2SPacketCommandLogOutUserReq
	// S2CPacketCommandLogOutUserRes is
	S2CPacketCommandLogOutUserRes

	// C2SPacketCommandLobbyMsgReq is
	C2SPacketCommandLobbyMsgReq
	// S2CPacketCommandLobbyMsgRes is
	S2CPacketCommandLobbyMsgRes

	// C2SPacketCommandRoomCreateReq is
	C2SPacketCommandRoomCreateReq
	// S2CPacketCommandRoomCreateRes is
	S2CPacketCommandRoomCreateRes

	// C2SPacketCommandRoomEnterReq is
	C2SPacketCommandRoomEnterReq
	// S2CPacketCommandRoomEnterRes is
	S2CPacketCommandRoomEnterRes

	// C2SPacketCommandRoomLeaveReq is
	C2SPacketCommandRoomLeaveReq
	// S2CPacketCommandRoomLeaveRes is
	S2CPacketCommandRoomLeaveRes

	// C2SPacketCommandRoomMsgReq is
	C2SPacketCommandRoomMsgReq
	// S2CPacketCommandRoomMsgRes is
	S2CPacketCommandRoomMsgRes

	// S2CPacketCommandSystemMsgSend is
	S2CPacketCommandSystemMsgSend
)

// C2SPCLobbySendMsgReq is  C2SPacketCommandLoginUserReq packet struct
type C2SPCLobbySendMsgReq struct {
	Msg string
}

// S2CPCLobbySendMsgRes is  S2CPacketCommandLoginUserRes packet struct
type S2CPCLobbySendMsgRes struct {
	Result uint32
	Userid string
	Msg    string
}

// C2SPCRoomCreateReq is  C2SPacketCommandRoomEnterReq packet struct
type C2SPCRoomCreateReq struct {
	RoomName string
}

// S2CPCRoomCreateRes is  C2SPacketCommandRoomLeaveReq packet struct
type S2CPCRoomCreateRes struct {
	Result      uint32
	RoomIdx     uint32
	RoomName    string
	EnterUserSn uint32
	EnterUserid string
}

// C2SPCRoomEnterReq is  C2SPacketCommandRoomEnterReq packet struct
type C2SPCRoomEnterReq struct {
	RoomName string
}

// S2CPCRoomEnterRes is  C2SPacketCommandRoomLeaveReq packet struct
type S2CPCRoomEnterRes struct {
	Result      uint32
	RoomIdx     uint32
	RoomName    string
	EnterUserSn uint32
	EnterUserid string
}

// C2SPCRoomLeaveReq is  C2SPacketCommandRoomEnterReq packet struct
type C2SPCRoomLeaveReq struct {
}

// S2CPCRoomLeaveRes is  S2CPacketCommandRoomLeaveRes packet struct
type S2CPCRoomLeaveRes struct {
	Result      uint32
	LeaveUserSn uint32
	LeaveUserid string
}

// C2SPCRoomSendMsgReq is  C2SPacketCommandRoomMsgReq packet struct
type C2SPCRoomSendMsgReq struct {
	RoomIdx uint32
	Msg     string
}

// S2CPCRoomSendMsgRes is  S2CPacketCommandRoomMsgRes packet struct
type S2CPCRoomSendMsgRes struct {
	Result uint32
	Userid string
	Msg    string
}

// S2CPCSystemMsgSend is  S2CPacketCommandSystemMsgSend packet struct
type S2CPCSystemMsgSend struct {
	Msg string
}
