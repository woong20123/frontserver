package examserverlogic

import (
	"example/examshare"
	"net"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// RegistServerLogic is regist Packet process logic from ChatServerMode
func RegistServerLogic(slm *tcpserver.ServerLogicManager) {
	RegistSystemServerLogic(slm)
}

// RegistSystemServerLogic is regist Packet process logic
func RegistSystemServerLogic(slm *tcpserver.ServerLogicManager) {
	// C2SPacketCommandLoginUserReq Packet Logic
	// 유저의 로그인 패킷 처리 작업 등록
	slm.RegistLogicfun(int32(examshare.Cmd_CS2FServerRegistRes), func(conn *net.TCPConn, p *packet.Packet) {
		res := examshare.CS2F_ServerRegistRes{}
		err := p.UnMarshalFromProto(&res)
		if err != nil {
			Logger().Println(err)
			return
		}
		return
	})
}
