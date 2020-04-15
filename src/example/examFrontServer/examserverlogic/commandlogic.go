package examserverlogic

import (
	"encoding/binary"
	"example/share"
	"log"
	"net"

	"github.com/woong20123/logicmanager"
	"github.com/woong20123/packet"
)

// RegistCommandLogic is
func RegistCommandLogic(lm *logicmanager.LogicManager) {
	lm.RegistLogicfun(share.C2SPacketCommandLoginUserReq, func(conn *net.TCPConn, p *packet.Packet) {
		var Userid string
		p.Read(binary.LittleEndian, &Userid)
		log.Println("C2SPacketCommandLoginUserReq userid = ", Userid)
		return
	})

	lm.RegistLogicfun(share.C2SPacketCommandGolobalSendMsgReq, func(conn *net.TCPConn, p *packet.Packet) {
		var creq share.C2SPCGolobalSendMsgReq
		binary.Read(p.GetByteBuf(), binary.LittleEndian, &creq)
		log.Println("C2SPacketCommandGolobalSendMsgReq msg = ", creq.Msg)
		return
	})
}
