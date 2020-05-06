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
