package examserverlogic

import (
	"example/examshare"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// RegistClientSessionLogic 은 Client에서 접속하는 세션에 대한 처리 로직 등록합니다.
func RegistClientSessionLogic(csessionhanlder *tcpserver.SessionHandler) {

	// 세션 연결시 ExamServer에서 해야 할 작업 등
	csessionhanlder.RegistConnectFunc(tcpserver.SessionStateEnum.OnConnected, func(s tcpserver.Session) {
		IssuedSn := Instance().ObjMgr().IssueSessionSn()
		Instance().ObjMgr().AddSession(IssuedSn, s)

		tcpclisession := s.(*tcpserver.TCPClientSession)
		if nil != tcpclisession {
			tcpclisession.SetSessionSn(IssuedSn)
		}
	})

	// 세션 종료시 ExamServer에서 해야 할 작업 등록
	csessionhanlder.RegistConnectFunc(tcpserver.SessionStateEnum.OnClosed, func(s tcpserver.Session) {
		eu := Instance().ObjMgr().FindUserByConn(s.Conn())
		if eu != nil {
			Instance().ObjMgr().DelUserByConn(s.Conn())
			eu.SetSession(nil)
			userID := eu.UserID()
			Instance().ObjMgr().DelUserString(&userID)
			Logger().Println("[", userID, "] 유저가 종료 하였습니다.")
		}
	})

	// 세션으로 데이터가 들어오면 해야 할 작업 등록
	csessionhanlder.RegistRecvFunc(func(s tcpserver.Session, buffer []byte, pos uint32) uint32 {
		var onPacket *packet.Packet = nil
		// 남은 버퍼에서 패킷을 조립할 수 있을 수도 있기 때문에 재호출
		for {
			pos, onPacket = packet.AssemblyFromBuffer(buffer, pos, uint32(examshare.Etc_ExamplePacketSerialkey))
			if onPacket == nil {
				break
			}
			tcpserver.Instance().ClientLogicManager().CallLogicFun(onPacket.Command(), s, onPacket)
		}
		return pos
	})

}

// RegistServerProxySessionLogic 은 ServerProxy 세션에 대한 처리 로직 등록합니다.
func RegistServerProxySessionLogic(spsessionhanlder *tcpserver.SessionHandler) {

	// 세션 연결시 ExamServer에서 해야 할 작업 등록
	spsessionhanlder.RegistConnectFunc(tcpserver.SessionStateEnum.OnConnected, func(s tcpserver.Session) {
		tcpsvrsession := s.(*tcpserver.TCPServerSession)
		if nil != tcpsvrsession {
			// ServerProxy Session Read goroutine start
			err := tcpserver.Instance().TCPServerMgr().AddTCPServerSession(tcpsvrsession)
			tcpserver.Instance().SendManager().RunSendToServerHandle(tcpsvrsession.Index())
			if err == nil {
				go tcpsvrsession.HandleRead()
			} else {
				Logger().Println(err.Error())
			}
		} else {
			println("[ChatServer] RegistConnectFunc fail")
			Logger().Println("[ChatServer] RegistConnectFunc fail")
		}
	})

	// 세션 종료시 ExamServer에서 해야 할 작업 등록
	spsessionhanlder.RegistConnectFunc(tcpserver.SessionStateEnum.OnClosed, func(s tcpserver.Session) {
	})

	// 세션으로 데이터가 들어오면 해야 할 작업 등록
	spsessionhanlder.RegistRecvFunc(func(s tcpserver.Session, buffer []byte, pos uint32) uint32 {
		var onPacket *packet.Packet = nil
		// 남은 버퍼에서 패킷을 조립할 수 있을 수도 있기 때문에 재호출
		for {
			pos, onPacket = packet.AssemblyFromBuffer(buffer, pos, uint32(examshare.Etc_ExamplePacketSerialkey))
			if onPacket == nil {
				break
			}
			tcpserver.Instance().ServerLogicManager().CallLogicFun(s.Index(), s, onPacket)
		}
		return pos
	})
}
