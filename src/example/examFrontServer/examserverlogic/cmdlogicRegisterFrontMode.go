package examserverlogic

import (
	"example/share"
	"net"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// FrontModeRegistCommandLogic is regist Packet process logic from ChatServerMode
func FrontModeRegistCommandLogic(lm *tcpserver.LogicManager) {
	frontModeRegistUserCommandLogic(lm)
	// ChatRoom 관련 패킷 로직 등록 함수
	frontModeRegistChatRoomCommandLogic(lm)
}

func frontModeRegistUserCommandLogic(lm *tcpserver.LogicManager) {
	// C2SPacketCommandLoginUserReq Packet Logic
	// 유저의 로그인 패킷 처리 작업 등록
	lm.RegistLogicfun(share.C2SPacketCommandLoginUserReq, func(conn *net.TCPConn, p *packet.Packet) {
		tcpserver.Instance().SendManager().SendToServerConn(share.TCPCliToSvrIdxChat, p)
		return
	})

	// C2SPacketCommandLobbyMsgReq Packet Logic
	// 로비에 전달하는 메시지 패킷 처리 작업 등록
	lm.RegistLogicfun(share.C2SPacketCommandLobbyMsgReq, func(conn *net.TCPConn, p *packet.Packet) {
		tcpserver.Instance().SendManager().SendToServerConn(share.TCPCliToSvrIdxChat, p)
		return
	})
}

func frontModeRegistChatRoomCommandLogic(lm *tcpserver.LogicManager) {

	// C2SPacketCommandRoomEnterReq Packet Logic =======================================================================
	// 유저의 방입장 패킷 처리 작업 등록 - 방이 없으면 생성합니다.
	lm.RegistLogicfun(share.C2SPacketCommandRoomEnterReq, func(conn *net.TCPConn, p *packet.Packet) {
		tcpserver.Instance().SendManager().SendToServerConn(share.TCPCliToSvrIdxChat, p)
		return
	})

	// C2SPacketCommandRoomCreateReq Packet Logic =======================================================================
	// 유저의 방 생성 패킷 처리 작업
	lm.RegistLogicfun(share.C2SPacketCommandRoomCreateReq, func(conn *net.TCPConn, p *packet.Packet) {
		tcpserver.Instance().SendManager().SendToServerConn(share.TCPCliToSvrIdxChat, p)
		return
	})

	// C2SPacketCommandRoomLeaveReq Packet Logic =======================================================================
	// 유저의 방 퇴장 패킷 처리 작업
	lm.RegistLogicfun(share.C2SPacketCommandRoomLeaveReq, func(conn *net.TCPConn, p *packet.Packet) {
		tcpserver.Instance().SendManager().SendToServerConn(share.TCPCliToSvrIdxChat, p)
		return
	})

	// C2SPacketCommandRoomMsgReq Packet Logic
	// 유저의 방에서 패킷을 전송합니다
	lm.RegistLogicfun(share.C2SPacketCommandRoomMsgReq, func(conn *net.TCPConn, p *packet.Packet) {
		tcpserver.Instance().SendManager().SendToServerConn(share.TCPCliToSvrIdxChat, p)
		return
	})
}
