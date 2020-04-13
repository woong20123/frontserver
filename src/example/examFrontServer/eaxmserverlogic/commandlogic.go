package eaxmserverlogic

import (
	"example/examlogic"

	"github.com/woong20123/logicmanager"
	"github.com/woong20123/packet"
)

func RegistCommandLogic(lm *logicmanager.LogicManager) {
	lm.RegistLogicfun(examlogic.C2SPacketCommandLoginUserReq, func(p *packet.Packet) {
		return
	})

	lm.RegistLogicfun(examlogic.C2SPacketCommandGolobalSendMsgReq, func(p *packet.Packet) {
		return
	})
}
