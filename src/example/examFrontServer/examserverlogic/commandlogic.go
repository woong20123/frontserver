package examserverlogic

import (
	"example/share"
	"log"
	"net"

	"github.com/woong20123/logicmanager"
	"github.com/woong20123/packet"
)

// RegistCommandLogic is
func RegistCommandLogic(lm *logicmanager.LogicManager) {
	lm.RegistLogicfun(share.C2SPacketCommandLoginUserReq, func(conn *net.TCPConn, p *packet.Packet) {
		log.Println("C2SPacketCommandLoginUserReq")
		return
	})

	lm.RegistLogicfun(share.C2SPacketCommandGolobalSendMsgReq, func(conn *net.TCPConn, p *packet.Packet) {
		log.Println("C2SPacketCommandGolobalSendMsgReq")
		return
	})
}
