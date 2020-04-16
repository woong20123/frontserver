package examclientlogic

import (
	"example/examClient/clientuser"
	"example/share"
	"log"
	"net"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// ContructLogicManager is
func ContructLogicManager(lm *tcpserver.LogicManager) {
	lm.RegistLogicfun(share.S2CPacketCommandLoginUserRes, func(conn *net.TCPConn, p *packet.Packet) {

		res := share.S2CPCLoginUserRes{}
		p.Read(&res.Result, &res.UserSn)
		res.UserID = p.ReadString()

		if res.Result == share.ResultSuccess {
			log.Println("Login Success sn = ", res.UserSn, " ID = ", res.UserID)
		} else {
			log.Println("Login Fail")
		}

		// 패킷을 확인하고 체크합니다.

		GetInstance().GetObjMgr().GetChanManager().SendChanUserState(clientuser.UserStateEnum.LoginSTATE, 200)
		return
	})

	lm.RegistLogicfun(share.S2CPacketCommandGolobalMsgRes, func(conn *net.TCPConn, p *packet.Packet) {
		log.Println("S2CPacketCommandGolobalSendMsgRes")
		return
	})

	lm.RunLogicHandle(1)
}
