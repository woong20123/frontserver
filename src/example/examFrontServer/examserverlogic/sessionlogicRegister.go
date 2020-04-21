package examserverlogic

import (
	"example/share"
	"net"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// RegistSessionLogic is
func RegistSessionLogic(sessionm *tcpserver.SessionMgr) {

	// 세션 연결시 ExamServer에서 해야 할 작업 등록
	sessionm.RegistConnectFunc(tcpserver.SessionStateEnum.OnConnected, func(conn *net.TCPConn) {
		eu := NewExamUser()
		eu.SetConn(conn)
		eu.SetState(UserStateEnum.ConnectedSTATE)
		GetInstance().GetObjMgr().AddUser(conn, eu)
	})

	// 세션 종료시 ExamServer에서 해야 할 작업 등록
	sessionm.RegistConnectFunc(tcpserver.SessionStateEnum.OnClosed, func(conn *net.TCPConn) {
		eu := GetInstance().GetObjMgr().FindUser(conn)
		if eu != nil {
			eu.SetConn(nil)
			userID := eu.GetUserID()
			GetInstance().GetObjMgr().DelUserString(&userID)
			GetLogger().Println("[", userID, "] 유저가 종료 하였습니다.")
		}
		GetInstance().GetObjMgr().DelUser(conn)
	})

	// 세션으로 데이터가 들어오면 해야 할 작업 등록
	sessionm.RegistRecvFunc(func(conn *net.TCPConn, buffer []byte, pos uint32) uint32 {
		var onPacket *packet.Packet = nil
		// 남은 버퍼에서 패킷을 조립할 수 있을 수도 있기 때문에 재호출
		for {
			pos, onPacket = packet.AssemblyFromBuffer(buffer, pos, share.ExamplePacketSerialkey)
			if onPacket == nil {
				break
			}
			tcpserver.Instance().LogicManager().CallLogicFun(onPacket.Command(), conn, onPacket)
		}
		return pos
	})
}
