package examserverlogic

import (
	"net"

	"github.com/woong20123/tcpserver"
)

// RegistSessionLogic is
func RegistSessionLogic(sessionm *tcpserver.SessionMgr) {

	// 세션 연결시 ExamServer에서 해야 할 작업 등록
	sessionm.RegistStateFunc(tcpserver.SessionStateEnum.OnConnected, func(conn *net.TCPConn) {

	})

	// 세션 종료시 ExamServer에서 해야 할 작업 등록
	sessionm.RegistStateFunc(tcpserver.SessionStateEnum.OnClosed, func(conn *net.TCPConn) {

	})
}
