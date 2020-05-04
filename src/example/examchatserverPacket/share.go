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
